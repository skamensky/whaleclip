export interface Entry{
    content : string,
    lastCopied : Number
    id:string
}


export interface EntryProps{
    entries: Entry[]
}