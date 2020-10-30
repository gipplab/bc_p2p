use itertools::Itertools;
use crate::sem_scholar_utils::doc::Reference;

pub fn create_k2_sets(bib: Vec<Reference>) -> Vec<Vec<String>> {
    //Extract IDs to Collection
    let ids: Vec<String> = bib.iter().map(|r| r.paper_id.clone()).collect();
    println!("Number of unique Refs: {}",ids.len());
    let sets_of_two = ids.into_iter().combinations(2).collect();

    return sets_of_two;
}