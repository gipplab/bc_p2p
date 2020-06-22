//! Type `PUT my-key my-value`
//! Type `GET my-key`
//! Type `BATCH my-path-to-csv`
//! Type `CHECK my-path-to-csv`
//! Close with Ctrl-c.

mod semscholar;

use async_std::{io, task};
use futures::prelude::*;
use libp2p::kad::record::store::MemoryStore;
use libp2p::kad::record::store::MemoryStoreConfig;
//no persistent DB
use libp2p::kad::{record::Key, Kademlia, KademliaEvent, PutRecordOk, Quorum, Record};
// flex DTH algo, record = key != hash, callback event, ...
use libp2p::{
    NetworkBehaviour,
    PeerId, //hash
    Swarm, //Channel for our cluster
    build_development_transport, //?
    identity, //PubKey
    mdns::{Mdns, MdnsEvent}, // initial peer discovery = only in LAN or VPN
    swarm::NetworkBehaviourEventProcess, // ?
};
use std::{error::Error, task::{Context, Poll}, fs};
use std::path::Path;
use chrono::Local;
//use std::borrow::Borrow;
//use std::sync::mpsc::Receiver;

// We create a custom network behaviour that combines Kademlia and mDNS.
#[derive(NetworkBehaviour)] // behaviour = interface
struct MyBehaviour {
    kademlia: Kademlia<MemoryStore>,
    // non persistent key value store
    mdns: Mdns, //peer discovery - not usable in a open setting - hard coded peers needed
}

impl NetworkBehaviourEventProcess<MdnsEvent> for MyBehaviour {
    // implements the Peer-discovery interface method of NetworkBehaviour
// Called when `mdns` produces an event.
    fn inject_event(&mut self, event: MdnsEvent) {
        if let MdnsEvent::Discovered(list) = event {
            for (peer_id, multiaddr) in list {  //Multiaddr works with a variant of addresses (IPv4 for TCP in our case)
                self.kademlia.add_address(&peer_id, multiaddr);
            }
        }
    }
}

impl NetworkBehaviourEventProcess<KademliaEvent> for MyBehaviour {
    // implements the DHT interface method of NetworkBehavior
// Called when `kademlia` produces an event.
    fn inject_event(&mut self, message: KademliaEvent) {
        match message {
            KademliaEvent::GetRecordResult(Ok(result)) => {
                for Record { key, value, publisher, .. } in result.records {
                    println!(
                        "Got record {:?} {:?} from Publisher {:?} at {:?}",
                        std::str::from_utf8(key.as_ref()),
                        std::str::from_utf8(value.as_ref()),
                        publisher,  //PRINT THE Publisher PEER
                        Local::now()
                    );
                }
            }
            KademliaEvent::GetRecordResult(Err(err)) => {
                eprintln!("Failed to get record: {:?}", err);
            }
            KademliaEvent::PutRecordResult(Ok(PutRecordOk { key })) => {
                println!(
                    "Successfully put record {:?} at {:?}",
                    std::str::from_utf8(key.as_ref()),
                    Local::now()
                );
            }
            KademliaEvent::PutRecordResult(Err(err)) => {
                eprintln!("Failed to put record: {:?}", err);
            }
            _ => {}
        }
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> { // return type "Result" for debug error handling
    let citations_result = semscholar::get_citations_of().await;
    println!("Lookup = {:?}", citations_result);

    env_logger::init(); //console logs

    // Create a random key for ourselves.
    let local_key = identity::Keypair::generate_ed25519();
    let local_peer_id = PeerId::from(local_key.public());

    // Set up a an encrypted DNS-enabled TCP Transport over the Mplex protocol.
    let transport = build_development_transport(local_key.clone())?; // Transport layer ist variable, for our use case of small keys TCP is not ideal


    // Create a swarm to manage peers and events.
    let mut swarm = {  // swarm is like a channel - channel initialized with own peer_id
        // Create a Kademlia behaviour.

        //libp2p_kad::record::store
        let mem_config = MemoryStoreConfig {
            max_records: 250000, // default is 1024 - with 250000 we can support 500 features on k2
            max_value_bytes: 65 * 1024,
            max_provided_keys: 1024,
            max_providers_per_key: 20, //Kademilia standard, could be smaller then 20 for the low peer count in the pilot phase
        };

        let store = MemoryStore::with_config(local_peer_id.clone(), mem_config );
        //let store = MemoryStore::new(local_peer_id.clone());


        let kademlia = Kademlia::new(local_peer_id.clone(), store); //keys and routing table
        let mdns = Mdns::new()?;
        let behaviour = MyBehaviour { kademlia, mdns };
        Swarm::new(transport, behaviour, local_peer_id.clone())
    };

    // Listen on all interfaces and whatever port the OS assigns.

    //TCP
    Swarm::listen_on(&mut swarm, "/ip4/0.0.0.0/tcp/0".parse()?)?; //listening for mdns results on addr

    // Save cli input - blocking when busy
    helper_safe_cli(&mut swarm, local_peer_id)
}

fn helper_safe_cli(swarm: &mut Swarm<MyBehaviour, PeerId>, local_peer_id: PeerId) -> Result<(), Box<dyn Error>> {
    // Read full lines from stdin
    let mut stdin = io::BufReader::new(io::stdin()).lines(); //cli input + buffer for increased performance
    let mut upload_buffer: Vec<Key> = Vec::new();

    let mut listening = false;
    task::block_on(future::poll_fn(move |cx: &mut Context| { //handle input as async task + context for timeout handling
        // blocking all new inputs until task is complete - poll fn checks for task completion
        loop {
            match stdin.try_poll_next_unpin(cx)? { //handle input, empty input, and pending input
                Poll::Ready(Some(line)) => handle_input_line(&mut swarm.kademlia, line, local_peer_id.clone(),  &mut upload_buffer), //execute input command
                Poll::Ready(None) => panic!("Stdin closed"),
                Poll::Pending => break
            }
        }
        loop {
            match swarm.poll_next_unpin(cx) { //blocking until saving task has finished
                Poll::Ready(Some(event)) => println!("{:?}", event),
                Poll::Ready(None) => return Poll::Ready(Ok(())),
                Poll::Pending => {
                    if !listening {
                        if let Some(a) = Swarm::listeners(&swarm).next() {
                            println!("Listening on {:?}", a);
                            listening = true;
                        }
                    }
                    break;
                }
            }
        }

        if upload_buffer.len() > 0 {
            println!("Upload Buffer: {}", upload_buffer.len());
            let value = Vec::from(local_peer_id.to_base58());
            let key = upload_buffer.pop().unwrap().to_owned();
            let record = Record {
                key,
                value,
                publisher: Some(local_peer_id.to_owned()), // USEFUL FOR TRACEABILITY AND SPAM-PROTECTION
                expires: None, //stays in memory for ever + periodic replication and republication - Date as std::time::Instant
            };
            swarm.kademlia.put_record(record, Quorum::One); // Quorum = min replication factor specifies the minimum number of distinct nodes that must be successfully contacted in order for a query to succeed.
        }
        Poll::Pending
    }))
}

// Handle commands
fn handle_input_line(kademlia: &mut Kademlia<MemoryStore>, line: String, local_peer_id: PeerId, upload_buffer: &mut Vec<Key>) {
    let mut args = line.split(" ");

    match args.next() {
        Some("CHECK") => { //TODO: Refactor duplicated iterator
            println!("Started batch CHECK at: {}",Local::now());
            let feature_file_path = {
                match args.next() {
                    Some(feature_file_path) => Path::new(feature_file_path),
                    None => {
                        println!("No PATH provided, please set a PATH e.g.: ../../data/1000_k1.csv");
                        return;
                    }
                }
            };
            println!("Got path: {}", &feature_file_path.display());

            // Read CSV + handle invalid path
            let features = fs::read_to_string(&feature_file_path);
            let features = match features {
                Ok(f) => f,
                Err(_) => {
                    // /Users/corihle/GIT/SwarmBC/49222.csv
                    println!("No valid PATH provided, please set a PATH e.g.: /Users/corihle/GIT/bc_p2p/data/1000_k1.csv");
                    return;
                }
            };

            // Handle the 2 column CSV input
            let v: Vec<&str> = features.split(',').collect();
            for s in v.to_owned() {
                if s != v.first().unwrap().to_owned(){ //filter first
                    let key = Key::new(&s.lines().next().unwrap()
                        .chars().filter(|c| !c.is_whitespace())
                        .collect::<String>()
                    );
                    kademlia.get_record(&key, Quorum::One);
                }
            }
        }
        Some("BATCH") => {
            println!("Started BATCH upload at: {}",Local::now());
            let feature_file_path = {
                match args.next() {
                    Some(feature_file_path) => Path::new(feature_file_path),
                    None => {
                        println!("No PATH provided, please set a PATH e.g.: /Users/corihle/GIT/bc_p2p/data/1000_k1.csv");
                        return;
                    }
                }
            };
            println!("Got path: {}", &feature_file_path.display());

            // Read CSV + handle invalid path
            let features = fs::read_to_string(&feature_file_path);
            let features = match features {
                Ok(f) => f,
                Err(_) => {
                    // /Users/corihle/GIT/SwarmBC/49222.csv
                    println!("No valid PATH provided, please set a PATH e.g.: /Users/corihle/GIT/bc_p2p/data/1000_k1.csv");
                    return;
                }
            };

            // Handle the 2 column CSV input
            let v: Vec<&str> = features.split(',').collect();
            for s in v.to_owned() {
                if s != v.first().unwrap().to_owned(){ //filter first
                    let key = Key::new(&s.lines().next().unwrap()
                        .chars().filter(|c| !c.is_whitespace())
                        .collect::<String>()
                    );
                    upload_buffer.push(key);
                }
            }
        }
    Some("GET") => {
        let key = {
            match args.next() {
                Some(key) => Key::new(&key),
                None => {
                    eprintln!("Expected key");
                    return;
                }
            }
        };
        kademlia.get_record(&key, Quorum::One);
    }
    Some("PUT") => {
        let key = {
            match args.next() {
                Some(key) => Key::new(&key),
                None => {
                    eprintln!("Expected key");
                    return;
                }
            }
        };
        let value = {
            match args.next() {
                Some(value) => value.as_bytes().to_vec(),
                None => {
                    eprintln!("Expected value");
                    return;
                }
            }
        };
        let record = Record {
            key,
            value,
            publisher: Some(local_peer_id), // USEFUL FOR TRACEABILITY AND SPAM-PROTECTION
            expires: None, //stays in memory for ever + periodic replication and republication
        };
        kademlia.put_record(record, Quorum::One); // Quorum = min replication factor specifies the minimum number of distinct nodes that must be successfully contacted in order for a query to succeed.
    }
    _ => {
    eprintln ! ("expected GET or PUT");
    }
}
}
