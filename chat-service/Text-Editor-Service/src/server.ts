import express from "express";
import dotenv from "dotenv";

import type { Request, Response } from "express";
import {WebSocketServer} from "ws";
import {WebSocket} from "ws"
import type {Appmessage } from "./type.js";
dotenv.config();

const app = express();
const PORT = process.env.PORT || 3000;

// Middleware
app.use(express.json());
const wss = new WebSocketServer({ port: 3005 });

// wss.on("connection", (socket:WebSocket) => {

// socket.on("message",(msg:Buffer)=>{
//    const msgStr = Buffer.from(JSON.stringify(msg));
// console.log("Received message:", msg);

// var temp = JSON.parse(msg.toString());
//     console.log("Parsed message:", temp);
// })
//  socket.on("close", () => {
//     console.log("Client disconnected");
//   });



// })


const  rooms: Map<string,Set<WebSocket>> = new Map();

  

  

wss.on("connection",(socket:WebSocket)=>{
    console.log("New client connected");
    socket.on("message",(msg:Buffer)=>{
      wss.clients.forEach((client:WebSocket)=>{
        if (client.readyState === WebSocket.OPEN) {
          client.send(msg);
        }
      })
    const msgStr = msg.toString();
const msgJson = JSON.parse(msgStr);

      console.log("Received message:", msgStr);
   try{
      
     if(msgJson.type === "join-room"){
      console.log("Join message received:", msgJson);
    const  room = msgJson.room
       if(!rooms.has(room)){
          rooms.set(room,new Set<WebSocket>());
 
       } if(rooms.has(room)!== null){
        rooms.get(room)!.add(socket)
       }
        
       console.log(`Client joined room: ${rooms}`);
       for (const [key,value] of rooms){
        console.log(`Room: ${key}, Clients: ${Array.from(value).length}`);
       }
         
         


        socket.send(JSON.stringify({ type: "joined", room }));
        
        
     }

      if(msgJson.type ==="message"){
          console.log("Message received:", msgJson);
         const roomId = msgJson.roomId;
         const message = msgJson.message;
         //get the clients in the room
         if(!rooms.has(roomId)){
            console.log(`Room ${roomId} does not exist.`);
            return;
         }
            const clinetIntheRoom = rooms.get(roomId)
             if(clinetIntheRoom){
               console.log(`Clients in room ${roomId}: ${Array.from(clinetIntheRoom).length}`);
                for (const client of clinetIntheRoom) {
                  if (client && client.readyState === WebSocket.OPEN) {
   client.send(JSON.stringify({ type: "message", roomId, message }));


             }

      }
   }
  }
  else{
    console.log("found no room:", msgJson);
  }
}
   catch(error){
   throw new Error("Invalid message format");
   }

    })
})


// Routes const room = (message as Message).room;
app.get("/", (req: Request, res: Response) => {
  res.send("Hello from Express + TypeScript!");
});

// Start server
app.listen(PORT, () => {
  console.log(`🚀 Server running at http://localhost:${PORT}`);
});
