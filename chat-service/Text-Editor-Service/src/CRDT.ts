import { get } from "http";

type ID = [agent: string, seq: number //like the user  device
]
type Item = {
    id: ID,
    content: string,
    originLeft: ID | null,
    originRight: ID | null,
    DELETED?: boolean
}
type Version = Record<string, number>;
//this Document obj
//contains all the data i have written
//like it handles the state of document
// in the Server
type Document = {
    content: Item[];
    version: Version;
};



function createDoc(): Document {
    return {
        content: [],
        version: {}
    }
}

function InsertCharOne( doc: Document,pos:number,  agent:string , text:string)  {
// by this intial seqence  for the first agent will be  0

//  const index = findIndex(doc,pos,true)
let index = 0
 const  seq =  (doc.version[agent]??-1)+1
   CombineTheInserts( doc, {
    id: [agent, seq],
    content:text,
    originLeft: doc.content[index-1]?.id??null,// if ?? null then add 1
    originRight: doc.content[index]?.id??null,
  })
    

}
  function insertChar(doc:Document,text:string,agent:string){
      const  content = [...text]
      console.log("InsertCharOne", content);
    let pos = 0
        for (const it of content){
            InsertCharOne( doc,pos, agent, it);
             pos++
        }

  }
 insertChar(createDoc(), "Hello", "agent1");

 function  CombineTheInserts( doc: Document,item:Item) {
 console.log("CombineTheInserts", item);
 }

function GetContent(doc: Document) {
    // for (const item of doc.content) {
    //     console.log(`ID: ${item.id}, Content: ${item.content}, OriginLeft: ${item.originLeft}, OriginRight: ${item.originRight}`);
    // }
  let content:string = ""
  
     doc.content.forEach(item =>{
         if(!item.DELETED){
              content = content+ item.content + " ";
  console.log(`ID: ${item.id}, Content: ${item.content}, OriginLeft: ${item.originLeft}, OriginRight: ${item.originRight}`);
         }
       

     })
   return content
}



// function Test() {
//     const item: Item = {
//         id: ["agent1", 1],
//         content: "Sample content",
//         originLeft: ["agent1", 0],
//         originRight: ["agent1", 2]
//     };
//     GetContent({content:[item], version: {agent1: 1}});
//     console.log(item);
// }
// Test();