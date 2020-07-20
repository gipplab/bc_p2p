use std::collections::HashSet;
use crate::semscholar;
use futures::future::ok;
use combinations::Combinations;
use std::cmp::Ordering;
use std::iter::FromIterator;
use futures::SinkExt;
use itertools::Itertools;

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Root {
    #[serde(rename = "abstract")]
    pub abstract_field: Option<String>,
    pub arxiv_id: Option<::serde_json::Value>,
    pub authors: Option<Vec<Author>>,
    pub citation_velocity: Option<i64>,
    pub citations: Option<Vec<Citation>>,
    pub corpus_id: Option<i64>,
    pub doi: Option<String>,
    pub fields_of_study: Option<Vec<String>>,
    pub influential_citation_count: Option<i64>,
    #[serde(rename = "is_open_access")]
    pub is_open_access: Option<bool>,
    #[serde(rename = "is_publisher_licensed")]
    pub is_publisher_licensed: Option<bool>,
    pub paper_id: String,
    pub references: Option<Vec<Reference>>,
    pub title: Option<String>,
    pub topics: Option<Vec<Topic>>,
    pub url: Option<String>,
    pub venue: Option<String>,
    pub year: Option<i64>,
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Author {
    pub author_id: Option<String>,
    pub name: Option<String>,
    pub url: Option<String>,
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Citation {
    pub arxiv_id: Option<::serde_json::Value>,
    pub authors: Option<Vec<Author2>>,
    pub doi: Option<String>,
    pub intent: Vec<String>,
    pub is_influential: Option<bool>,
    pub paper_id: String,
    pub title: Option<String>,
    pub url: Option<String>,
    pub venue: Option<String>,
    pub year: Option<i64>,
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Author2 {
    pub author_id: Option<String>,
    pub name: Option<String>,
    pub url: Option<String>,
}

#[derive(Default, Debug, Clone, PartialEq, Ord, Eq, PartialOrd, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Reference {
    pub arxiv_id: Option<String>,
    pub authors: Option<Vec<Author3>>,
    pub doi: Option<String>,
    pub intent: Option<Vec<String>>,
    pub is_influential: Option<bool>,
    pub paper_id: String,
    pub title: Option<String>,
    pub url: Option<String>,
    pub venue: Option<String>,
    pub year: Option<i64>,
}

#[derive(Default, Debug, Clone, PartialEq, Ord, Eq, PartialOrd, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Author3 {
    pub author_id: Option<String>,
    pub name: Option<String>,
    pub url: Option<String>,
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Topic {
    pub topic: Option<String>,
    pub topic_id: Option<String>,
    pub url: Option<String>,
}

const API_URL: &'static str = "https://api.semanticscholar.org/v1/paper/";

pub async fn get_id_from_doi(doi: &str) -> Result<String, anyhow::Error> {
    let doc_url = format!("{}{}", API_URL, doi);
    let body = reqwest::get(&doc_url)
        .await?
        .text()
        .await?;

    let doc: Root = serde_json::from_str(body.as_ref())?;
    let id: String = doc.paper_id;

    Ok(id)
}

pub async fn get_all_references_by_id(id: &str) -> Result<Vec<Reference>, anyhow::Error> {
    let doc_url = format!("{}{}", API_URL, id);
    let body = reqwest::get(&doc_url)
        .await?
        .text()
        .await?;
    let doc: Root = serde_json::from_str(body.as_ref())?;
    Ok(doc.references.unwrap())
}

pub fn create_k2_sets(bib: Vec<Reference>) -> Vec<Vec<String>> {
    //Extract IDs to Collection
    let ids: Vec<String> = bib.iter().map(|r| r.paper_id.clone()).collect();
    let sets_of_two = ids.into_iter().combinations(2).collect();

    return sets_of_two
}

pub async fn get_all_citations_by_reference_id(id: &str) -> Result<Vec<Citation>, anyhow::Error> {
    //TODO: handle connections over a reqwest::Client for single handshake
    let doc_url = format!("{}{}", API_URL, id);
    let body = reqwest::get(&doc_url)
        .await?
        .text()
        .await?;
    let doc: Root = serde_json::from_str(body.as_ref())?;
    Ok(doc.citations.unwrap())
}

// pub async fn public_bc() {
//     // ---------------------------------------
//     // PUBLIC CENTRALIZED - REFERENCE - STORE
//     // ---------------------------------------
//
//     // Test API reachability by getting a paper id
//     let id = get_id_from_doi("10.1016/j.jnca.2020.102630").await?;
//     println!("Got Paper ID: {}", id);
//
//     // Get all references of a doc
//     let doc_refs: Vec<semscholar::Reference> = get_all_references_by_id(&*id).await?;
//
//     // Extract Paper_IDs
//     let mut doc_ref_ids: HashSet<String> = HashSet::new();
//     for r in &doc_refs {
//         doc_ref_ids.insert(r.paper_id.clone());
//     }
//
//     // Iterate through all references
//     for r in &doc_refs {
//         // Get all docs which cite a reference (BC > 0)
//         let citations_ref = get_all_citations_by_reference_id(&*r.paper_id).await?;
//         //TODO: error handling
//
//         // Get all references of the co-citing docs (Check BC)
//         for co_r in &citations_ref {
//             println!("Checking {}...", co_r.paper_id);
//
//             // Todo: retry
//
//             // Check if the paper_id is still accessible
//             let co_cite_refs = get_all_references_by_id(&*co_r.paper_id).await;
//             match &co_cite_refs {
//                 Ok(refs)=> {
//                     ();
//                 },
//                 Err(e)=> {
//                     println!("Citation not found");
//                     continue;
//                 }
//                 _ => {}
//             }
//
//             // Extract Paper_IDs
//             let mut co_doc_ref_ids = HashSet::new();
//             for r in co_cite_refs? {
//                 co_doc_ref_ids.insert(r.paper_id);
//             }
//
//             // Find intersection between co-citing paper and source paper
//             let intersection_set: HashSet<_> = doc_ref_ids.intersection(&co_doc_ref_ids).collect();
//
//         }
//     }
//
//     // let citations_result = semscholar::get_citations_of().await;
//     //println!("Lookup = {:?}", citations_result);
//     return bc
// }

