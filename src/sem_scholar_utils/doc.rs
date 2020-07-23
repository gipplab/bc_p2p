#[derive(
    Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize,
)]
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

#[derive(
    Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize,
)]
#[serde(rename_all = "camelCase")]
pub struct Author {
    pub author_id: Option<String>,
    pub name: Option<String>,
    pub url: Option<String>,
}

#[derive(
    Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize,
)]
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

#[derive(
    Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize,
)]
#[serde(rename_all = "camelCase")]
pub struct Author2 {
    pub author_id: Option<String>,
    pub name: Option<String>,
    pub url: Option<String>,
}

#[derive(
    Default,
    Debug,
    Clone,
    PartialEq,
    Ord,
    Eq,
    PartialOrd,
    serde_derive::Serialize,
    serde_derive::Deserialize,
)]
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

#[derive(
    Default,
    Debug,
    Clone,
    PartialEq,
    Ord,
    Eq,
    PartialOrd,
    serde_derive::Serialize,
    serde_derive::Deserialize,
)]
#[serde(rename_all = "camelCase")]
pub struct Author3 {
    pub author_id: Option<String>,
    pub name: Option<String>,
    pub url: Option<String>,
}

#[derive(
    Default, Debug, Clone, PartialEq, serde_derive::Serialize, serde_derive::Deserialize,
)]
#[serde(rename_all = "camelCase")]
pub struct Topic {
    pub topic: Option<String>,
    pub topic_id: Option<String>,
    pub url: Option<String>,
}