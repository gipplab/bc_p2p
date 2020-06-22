use serde:: { Deserialize };
use reqwest::Error;

#[derive(Deserialize,Debug)]
struct Doc{
    meta_data:String,
    doi:String
}

pub async fn get_citations_of() -> Result<(), Error>{

    let body = reqwest::get("https://api.semanticscholar.org/v1/paper/10.1038/nrn3241")
        .await?
        .text()
        .await?;

    println!("body = {:?}", body);

    Ok(())
}

