use chrono::Local;

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Root {
    #[serde(rename = "abstract")]
    pub abstract_field: String,
    pub arxiv_id: ::serde_json::Value,
    pub authors: Vec<Author>,
    pub citation_velocity: i64,
    pub citations: Vec<Citation>,
    pub corpus_id: i64,
    pub doi: String,
    pub fields_of_study: Vec<String>,
    pub influential_citation_count: i64,
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
    pub year: i64,
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
    pub is_influential: bool,
    pub paper_id: String,
    pub title: String,
    pub url: String,
    pub venue: String,
    pub year: i64,
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Author2 {
    pub author_id: Option<String>,
    pub name: String,
    pub url: Option<String>,
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Reference {
    pub arxiv_id: Option<String>,
    pub authors: Vec<Author3>,
    pub doi: Option<String>,
    pub intent: Vec<String>,
    pub is_influential: bool,
    pub paper_id: String,
    pub title: String,
    pub url: String,
    pub venue: String,
    pub year: i64,
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Author3 {
    pub author_id: Option<String>,
    pub name: String,
    pub url: Option<String>,
}

#[derive(Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Topic {
    pub topic: String,
    pub topic_id: String,
    pub url: String,
}

pub async fn get_citations_of() -> Result<(), anyhow::Error> {
    let test_id: String = "fdd81cddbb086f377d3640581bcec58ec8f22c61".parse()?;


    // test for 100 queries
    let mut next_id: String = test_id.to_owned();
    println!("Start time {}", Local::now());
    for n in 0..100 {
        let doc_url = format!("https://api.semanticscholar.org/v1/paper/{}", next_id);

        let body = reqwest::get(&doc_url)
            .await?
            .text()
            .await?;

        let curr_id: String = String::from(next_id.clone());
        let doc: Root = serde_json::from_str(body.as_ref())?;
        let citations = doc.citations;
        next_id = citations[0].paper_id.to_owned();
        //let current_id = next_id.clone();

        let citations_count = citations.len();

        // Access parts of the data by indexing with square brackets.
        println!("{} papers cited the document {}. Query number {}", citations_count, curr_id, n);
    }
    println!("End time {}", Local::now());

    Ok(())
}

pub async fn get_id_from_doi(doi: &str) -> Result<String, anyhow::Error> {
    let doc_url = format!("https://api.semanticscholar.org/v1/paper/{}", doi);
    let body = reqwest::get(&doc_url)
        .await?
        .text()
        .await?;

    let doc: Root = serde_json::from_str(body.as_ref())?;
    let id: String = doc.paper_id;

    Ok(id)
}

pub async fn get_all_references_by_id() -> Result<(), anyhow::Error> {
    Ok(())
}


