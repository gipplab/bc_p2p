
#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Root {
    #[serde(rename = "abstract")]
    pub abstract_field: String,
    pub arxiv_id: ::serde_json::Value,
    pub authors: Vec<Author>,
    pub citation_velocity: Option<i64>,
    pub citations: Vec<Citation>,
    pub corpus_id: Option<i64>,
    pub doi: Option<String>,
    pub fields_of_study: Vec<String>,
    pub influential_citation_count: Option<i64>,
    #[serde(rename = "is_open_access")]
    pub is_open_access: bool,
    #[serde(rename = "is_publisher_licensed")]
    pub is_publisher_licensed: bool,
    pub paper_id: String,
    pub references: Vec<Reference>,
    pub title: String,
    pub topics: Vec<Topic>,
    pub url: String,
    pub venue: String,
    pub year: Option<i64>,
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Author {
    pub author_id: String,
    pub name: String,
    pub url: String,
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Citation {
    pub arxiv_id: ::serde_json::Value,
    pub authors: Vec<Author2>,
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

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Reference {
    pub arxiv_id: Option<String>,
    pub authors: Vec<Author3>,
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
    Ok(doc.references)
}

pub async fn get_all_citations_by_reference_id(id: &str) -> Result<Vec<Citation>, anyhow::Error> {
    let doc_url = format!("{}{}", API_URL, id);
    let body = reqwest::get(&doc_url)
        .await?
        .text()
        .await?;
    let doc: Root = serde_json::from_str(body.as_ref())?;
    Ok(doc.citations)
}


