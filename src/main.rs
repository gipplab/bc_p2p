//! bc_p2p
//! ======
//!
//! A peer-to-peer similarity detection experiment
//!
//! CLI Usage:
//! Type `PUT my-key my-value`
//! Type `GET my-key`
//! Type `BATCH my-path-to-csv`
//! Type `CHECK my-path-to-csv`
//! Close with Ctrl-c.

mod sem_scholar_utils;
use crate::sem_scholar_utils::api::{get_all_references_by_id, get_all_citations_by_reference_id};
use crate::sem_scholar_utils::doc::Reference;
use crate::sem_scholar_utils::sbc::create_k2_sets;

use async_std::{io, task};
use futures::prelude::*;
use libp2p::kad::record::store::MemoryStore;
use libp2p::kad::record::store::MemoryStoreConfig;
use libp2p::kad::{
    Kademlia,
    KademliaEvent,
    PeerRecord,
    PutRecordOk,
    QueryResult,
    Quorum,
    Record,
    record::Key,
};
use libp2p::{
    NetworkBehaviour,
    PeerId,
    Swarm,
    build_development_transport,
    identity,
    mdns::{Mdns, MdnsEvent},
    swarm::NetworkBehaviourEventProcess
};
use std::{error::Error, task::{Context, Poll}, fs, env};
use std::collections::{HashMap, HashSet};
use sha1::Sha1;
use sha1::Digest;
use std::io::Read;
use std::iter::FromIterator;
use chrono::{Local};
use chrono::DateTime;
use std::path::Path;
use lazy_static::lazy_static;
use mut_static::MutStatic;

pub struct Timer {
    startTime : DateTime<Local>
}

impl Timer {
    pub fn new(t: DateTime<Local>) -> Self{
        Timer { startTime: t }
    }
    pub fn getvalue(&self) -> DateTime<Local> { self.startTime }
    pub fn setvalue(&mut self, t: DateTime<Local>) { self.startTime = t }
}

lazy_static! {
    static ref MY_TIMER: MutStatic<Timer> = MutStatic::new();
}

// We create a custom network behaviour that combines Kademlia and mDNS.
#[derive(NetworkBehaviour)]
struct MyBehaviour {
    kademlia: Kademlia<MemoryStore>,
    mdns: Mdns,
}

impl NetworkBehaviourEventProcess<MdnsEvent> for MyBehaviour {
    // Called when `mdns` produces an event.
    fn inject_event(&mut self, event: MdnsEvent) {
        if let MdnsEvent::Discovered(list) = event {
            for (peer_id, multiaddr) in list {
                self.kademlia.add_address(&peer_id, multiaddr);
            }
        }
    }
}

impl NetworkBehaviourEventProcess<KademliaEvent> for MyBehaviour {
    // Called when `kademlia` produces an event.
    fn inject_event(&mut self, message: KademliaEvent) {
        match message {
            KademliaEvent::QueryResult { result, .. } => match result {
                QueryResult::GetRecord(Ok(ok)) => {
                    for PeerRecord { record: Record { key, value, publisher, .. }, ..} in ok.records {
                        println!(
                            "Got record {:?} {:?} from Publisher {:?} at {:?}",
                            std::str::from_utf8(key.as_ref()),
                            std::str::from_utf8(value.as_ref()),
                            publisher,  //PRINT THE Publisher PEER
                            Local::now() - MY_TIMER.read().unwrap().getvalue()
                        );
                    }
                }
                QueryResult::GetRecord(Err(err)) => {
                    eprintln!("Failed to get record: {:?}", err);
                }
                QueryResult::PutRecord(Ok(PutRecordOk { key })) => {
                    println!(
                        "Successfully put record {:?} at {:?}",
                        std::str::from_utf8(key.as_ref()),
                        Local::now() - MY_TIMER.read().unwrap().getvalue()
                    );
                }
                QueryResult::PutRecord(Err(err)) => {
                    eprintln!("Failed to put record: {:?}", err);
                }
                _ => {}
            }
            _ => {}
        }
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let args: Vec<String> = env::args().collect();
    env_logger::init();

    // Shared timer init
    MY_TIMER.set(Timer::new(Local::now())).unwrap();

    // --------------------
    // FILTERING
    // --------------------
    let k2_hashes;
    if args.len() > 1 {
        match args[1].as_str() {
            "filter" => {
                println!("filtering by document");
                k2_hashes = referenceFiltering().await;
                // TODO: add document ID input dialog
            },
            "json" => {
                println!("check public database with DOIs from JSON file");

               // paper_ids = get_paper_ids_from_json(arXid_json: &str)
                // TODO: add document ID input dialog
            },
            _ => println!("Standard Mode"), 
        }
    }

    // I need a list of dois

    // TODO: Check and upload k2_hashes to DHT
    // TODO:  originality ratio can be calculatedh as a numeric value

    // --------------------
    // P2P Upload
    // --------------------

    // Create a random key for ourselves.
    let local_key = identity::Keypair::generate_ed25519();
    let local_peer_id = PeerId::from(local_key.public());

    // Set up a an encrypted DNS-enabled TCP Transport over the Mplex protocol.
    let transport = build_development_transport(local_key)?;

    // Create a swarm to manage peers and events.
    let mut swarm = {
        // libp2p_kad::record::store
        let mem_config = MemoryStoreConfig {
            max_records: 250000, // default is 1024 - with 250000 we can support 500 features on k2
            max_value_bytes: 65 * 1024,
            max_provided_keys: 1024, // why not 20k?
            max_providers_per_key: 2, //Kademlia standard, could be smaller then 20 for the low peer count in the pilot phase
        };

        // Create a Kademlia behaviour.
        let store = MemoryStore::with_config(local_peer_id.clone(), mem_config );

        let kademlia = Kademlia::new(local_peer_id.clone(), store);
        let mdns = Mdns::new()?;
        let behaviour = MyBehaviour { kademlia, mdns };
        Swarm::new(transport, behaviour, local_peer_id.clone())
    };

    // Listen on all interfaces and whatever port the OS assigns.
    Swarm::listen_on(&mut swarm, "/ip4/0.0.0.0/tcp/0".parse()?)?;


    // SAVE CLI INPUT
    helper_safe_cli(&mut swarm, local_peer_id)
}

async fn referenceFiltering() -> Vec<String> {
// --------------------
    // REFERENCE FILTERING
    // --------------------

    // Get and print combinations from non-unique document
    let my_refs = get_all_references_by_id("97efafdb4a3942ab3efba53ded7413199f79c054").await.unwrap();


    // No new tuples in "Mathematical Formulae in Wikimedia Projects"
    //let my_refs = get_all_references_by_id("10.1145/3383583.3398557").await?.clone();
    //println!("Test document: {}", "10.1145/3383583.3398557");

    // No new tuples in "A First Step Towards Content Protecting Plagiarism Detection"
    // let my_refs = get_all_references_by_id("10.1145/3383583.3398620").await?.clone();
    // println!("Test document: {}", "10.1145/3383583.3398620");


    return filter_pub_refs(my_refs).await;
}

fn helper_safe_cli(swarm: &mut Swarm<MyBehaviour, PeerId>, local_peer_id: PeerId) -> Result<(), Box<dyn Error>> {
// Read full lines from stdin
    let mut stdin = io::BufReader::new(io::stdin()).lines();

    let mut upload_buffer: Vec<Key> = Vec::new();
    let mut listening = false;

    task::block_on(future::poll_fn(move |cx: &mut Context| {
        loop {
            match stdin.try_poll_next_unpin(cx)? {
                Poll::Ready(Some(line)) => handle_input_line(&mut swarm.kademlia, line, local_peer_id.clone(), &mut upload_buffer),
                Poll::Ready(None) => panic!("Stdin closed"),
                Poll::Pending => break
            }
        }
        loop {
            match swarm.poll_next_unpin(cx) {
                Poll::Ready(Some(event)) => println!("{:?}", event),
                Poll::Ready(None) => return Poll::Ready(Ok(())),
                Poll::Pending => {
                    if !listening {
                        if let Some(a) = Swarm::listeners(&swarm).next() {
                            println!("Listening on {:?}", a);
                            listening = true;
                        }
                    }
                    break
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

async fn filter_pub_refs(my_refs: Vec<Reference>) -> Vec<String> {
    let mut co_cite_refs_map: HashMap<String, HashSet<String>> = HashMap::new(); //reference; vector<paper_id>

    // Get all citations from own references
    for my_ref in my_refs.clone() {
        println!("Checking Reference: {}", my_ref.paper_id);
        let mut cits;
        match get_all_citations_by_reference_id(&*my_ref.paper_id).await {
            Ok(c_new) => cits = c_new,
            _ => cits = vec![],
        }

        let mut current_cits: HashSet<String>;
        match co_cite_refs_map.get(&*my_ref.paper_id) {
            Some(c_vec) => current_cits = c_vec.clone(),
            _ => current_cits = HashSet::new(),
        }

        // fill co-cite map
        for c in cits {
            current_cits.insert(c.paper_id);
        };
        co_cite_refs_map.insert(my_ref.paper_id, current_cits.clone());
    }
    // create k2 set
    let k2_sets = create_k2_sets(my_refs.clone());

    let mut k2_hashes: Vec<String> = vec![];

    // check entries for each k2 pair -> is it cited by the same doc_id
    for r in k2_sets {
        //println!("{} + {}", r[0], r[1]);

        // hash_a - id1,id2,id3,id4
        let all_docs_that_cite_a: HashSet<String> = co_cite_refs_map.get(&*r[0]).unwrap().clone();
        // hash_b - id1,id5,id6,id7
        let all_docs_that_cite_b: HashSet<String> = co_cite_refs_map.get(&*r[1]).unwrap().clone();

        let matching: HashSet<_> = all_docs_that_cite_a.intersection(&all_docs_that_cite_b).into_iter().clone().collect();

        // if not -> push to p2p hash table
        if matching.len() == 0 {
            println!("Publish private k2 set {} + {}", r[0], r[1]);
            let mut hash_output = Sha1::new().
                chain(&r[0]).
                chain(&r[1]).
                finalize();
            let hash = format!("{:x}", hash_output);
            println!("{}", hash);
            k2_hashes.push(hash);
        }
    }

    return k2_hashes;
}

// Handle commands
fn handle_input_line(kademlia: &mut Kademlia<MemoryStore>, line: String, local_peer_id: PeerId, upload_buffer: &mut Vec<Key>) {
    let mut args = line.split(" ");

    match args.next() {
        Some("CHECK") => { //TODO: Refactor duplicated iterator
            {
                let mut mut_timer = MY_TIMER.write().unwrap();
                mut_timer.setvalue(Local::now());
            }
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
            {
            let mut mut_timer = MY_TIMER.write().unwrap();
            mut_timer.setvalue(Local::now());
            }
            println!("Started BATCH upload at: {}",MY_TIMER.read().unwrap().getvalue());

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
            {
                let mut mut_timer = MY_TIMER.write().unwrap();
                mut_timer.setvalue(Local::now());
            }
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
            {
                let mut mut_timer = MY_TIMER.write().unwrap();
                mut_timer.setvalue(Local::now());
            }
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
            eprintln ! ("expected GET, PUT, BATCH, or CHECK");
        }
    }
}