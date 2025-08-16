import { get } from "http";

type ID = [agent: string, seq: number //like the user  device
]
type Item = {
    id: ID,
    content: string,
    originLeft: ID,
    originRight: ID,
    DELETED?: boolean
}
type Version = Record<string, number>;
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

 function GetContent(doc :Document) :Document {
      for  (const item of doc.content){
      console.log(`ID: ${item.id}, Content: ${item.content}, OriginLeft: ${item.originLeft}, OriginRight: ${item.originRight}`);
      }
      return doc
 }
  function InsertItem( doc:Document,item:Item):Document{

     
 return doc
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