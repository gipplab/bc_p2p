use crate::sem_scholar_utils::doc::*;

const API_URL: &'static str = "https://api.semanticscholar.org/v1/paper/";

pub async fn get_id_from_doi(doi: &str) -> Result<String, anyhow::Error> {
    let doc_url = format!("{}{}", API_URL, doi);
    let body = reqwest::get(&doc_url).await?.text().await?;

    let doc: Root = serde_json::from_str(body.as_ref())?;
    let id: String = doc.paper_id;

    Ok(id)
}

pub async fn get_all_references_by_id(id: &str) -> Result<Vec<Reference>, anyhow::Error> {
    let doc_url = format!("{}{}", API_URL, id);
    let body = reqwest::get(&doc_url).await?.text().await?;
    let doc: Root = serde_json::from_str(body.as_ref())?;
    Ok(doc.references.unwrap())
}

pub async fn get_all_citations_by_reference_id(
    id: &str,
) -> Result<Vec<Citation>, anyhow::Error> {
    //TODO: handle connections over a reqwest::Client for single handshake
    let doc_url = format!("{}{}", API_URL, id);
    let body = reqwest::get(&doc_url).await?.text().await?;
    let doc: Root = serde_json::from_str(body.as_ref())?;
    Ok(doc.citations.unwrap())
}