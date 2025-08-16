

export type OperationMessage = {
  type: "op";
  op: any; // later replace `any` with your operation structure
};
export type Message = {
  type: string;
  message?: any;
  roomId?: string;
  [key: string]: any;
};

export  type Appmessage = OperationMessage | Message;