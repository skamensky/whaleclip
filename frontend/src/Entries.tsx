import {Fragment, useState} from 'react';
import {Entry,EntryProps} from './types'

var t: {[i:string]:string}

function Entries(props:EntryProps){
    const entries =  props.entries.map((entry,idx)=>{
        return <p key={idx}>
            {entry.content}
        </p>
    })
    return <Fragment>{entries}</Fragment>
}

export default Entries