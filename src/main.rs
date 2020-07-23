
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
use std::{error::Error, task::{Context, Poll}};
use std::collections::{HashMap, HashSet};

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    env_logger::init();


// ---------------------------------------
    // REFERENCE FILTERING
    // ---------------------------------------

    // Get and print combinations from non-unique document
    //let my_refs = get_all_references_by_id("97efafdb4a3942ab3efba53ded7413199f79c054").await?.clone();
    //println!("Test document: {}", "97efafdb4a3942ab3efba53ded7413199f79c054");

    // No new tuples in "Mathematical Formulae in Wikimedia Projects"
    let my_refs = get_all_references_by_id("10.1145/3383583.3398557").await?.clone();
    println!("Test document: {}", "10.1145/3383583.3398557");

    // No new tuples in "A First Step Towards Content Protecting Plagiarism Detection"
    // let my_refs = get_all_references_by_id("10.1145/3383583.3398620").await?.clone();
    // println!("Test document: {}", "10.1145/3383583.3398620");

    filter_pub_refs(my_refs).await;

    //================================



    // Create a random key for ourselves.
    let local_key = identity::Keypair::generate_ed25519();
    let local_peer_id = PeerId::from(local_key.public());

    // Set up a an encrypted DNS-enabled TCP Transport over the Mplex protocol.
    let transport = build_development_transport(local_key)?;

    // We create a custom network behaviour that combines Kademlia and mDNS.
    #[derive(NetworkBehaviour)]
    struct MyBehaviour {
        kademlia: Kademlia<MemoryStore>,
        mdns: Mdns
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
                        for PeerRecord { record: Record { key, value, .. }, ..} in ok.records {
                            println!(
                                "Got record {:?} {:?}",
                                std::str::from_utf8(key.as_ref()).unwrap(),
                                std::str::from_utf8(&value).unwrap(),
                            );
                        }
                    }
                    QueryResult::GetRecord(Err(err)) => {
                        eprintln!("Failed to get record: {:?}", err);
                    }
                    QueryResult::PutRecord(Ok(PutRecordOk { key })) => {
                        println!(
                            "Successfully put record {:?}",
                            std::str::from_utf8(key.as_ref()).unwrap()
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

    // Create a swarm to manage peers and events.
    let mut swarm = {
        // Create a Kademlia behaviour.
        let store = MemoryStore::new(local_peer_id.clone());
        let kademlia = Kademlia::new(local_peer_id.clone(), store);
        let mdns = Mdns::new()?;
        let behaviour = MyBehaviour { kademlia, mdns };
        Swarm::new(transport, behaviour, local_peer_id)
    };

    // Read full lines from stdin
    let mut stdin = io::BufReader::new(io::stdin()).lines();

    // Listen on all interfaces and whatever port the OS assigns.
    Swarm::listen_on(&mut swarm, "/ip4/0.0.0.0/tcp/0".parse()?)?;

    // Kick it off.
    let mut listening = false;
    task::block_on(future::poll_fn(move |cx: &mut Context| {
        loop {
            match stdin.try_poll_next_unpin(cx)? {
                Poll::Ready(Some(line)) => handle_input_line(&mut swarm.kademlia, line),
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
        Poll::Pending
    }))
}

async fn filter_pub_refs(my_refs: Vec<Reference>) {
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
    // TODO: pre-filter set for
    let k2_sets = create_k2_sets(my_refs.clone());

    // check entries for each k2 pair -> is it cited by the same doc_id
    for r in k2_sets {
        //println!("{} + {}", r[0], r[1]);

        let a: HashSet<String> = co_cite_refs_map.get(&*r[0]).unwrap().clone();
        let b: HashSet<String> = co_cite_refs_map.get(&*r[1]).unwrap().clone();

        let matching: HashSet<_> = a.intersection(&b).into_iter().clone().collect();

        // if not -> push to p2p hash table
        if matching.len() == 0 {
            println!("Publish private k2 set {} + {}", r[0], r[1])
        }
    }
}

fn handle_input_line(kademlia: &mut Kademlia<MemoryStore>, line: String) {
    let mut args = line.split(" ");

    match args.next() {
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
                publisher: None,
                expires: None,
            };
            kademlia.put_record(record, Quorum::One).expect("Failed to store record locally.");
        }
        _ => {
            eprintln!("expected GET or PUT");
        }
    }
}