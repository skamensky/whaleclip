import {useState,useEffect} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {EventsOn,EventsOff} from "../wailsjs/runtime"
import Entries from "./Entries";
import {Entry,EntryProps} from './types'

function App() {

    const [entries,setEntries] = useState<Entry[]>([])

    useEffect(()=>{
        console.log("subscribing to newclip event")
        EventsOn("newClips",(data:EntryProps)=>{
            if(!data||!Object.keys(data).length){
                console.error(`No data in newclip event data received`,data)
                return
            }

            const newIds = new Set(data.entries.map(e=>e.id))
            setEntries(currentEntries=>{
                const oldEntriesDeduped = currentEntries.filter(e=>{
                    return !newIds.has(e.id)
                })
                const combinedEntries = [...oldEntriesDeduped,...data.entries]
                console.log(combinedEntries)

                return combinedEntries.sort((a,b)=>{
                    return Number(a.lastCopied>b.lastCopied)
                });
            })
        })

        return ()=>{
            console.log("unsubscribing from newclip event")
            EventsOff("newClips")
        }
    },[])

    return (
        <div id="App">
            <Entries entries={entries}></Entries>
        </div>
    )
}

export default App
