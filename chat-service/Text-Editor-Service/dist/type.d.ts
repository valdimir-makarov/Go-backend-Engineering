export type JoinMessage = {
    type: "join";
    room: string;
};
export type OperationMessage = {
    type: "op";
    op: any;
};
export type Message = JoinMessage | OperationMessage;
//# sourceMappingURL=type.d.ts.map