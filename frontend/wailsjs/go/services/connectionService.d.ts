// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {types} from '../models';
import {context} from '../models';

export function CreateGroup(arg1:string):Promise<types.JSResp>;

export function DeleteConnection(arg1:string):Promise<types.JSResp>;

export function DeleteGroup(arg1:string,arg2:boolean):Promise<types.JSResp>;

export function ExportConnections():Promise<types.JSResp>;

export function GetConnection(arg1:string):Promise<types.JSResp>;

export function ImportConnections():Promise<types.JSResp>;

export function ListConnection():Promise<types.JSResp>;

export function ListSentinelMasters(arg1:types.ConnectionConfig):Promise<types.JSResp>;

export function ParseConnectURL(arg1:string):Promise<types.JSResp>;

export function RenameGroup(arg1:string,arg2:string):Promise<types.JSResp>;

export function SaveConnection(arg1:string,arg2:types.ConnectionConfig):Promise<types.JSResp>;

export function SaveLastDB(arg1:string,arg2:number):Promise<types.JSResp>;

export function SaveRefreshInterval(arg1:string,arg2:number):Promise<types.JSResp>;

export function SaveSortedConnection(arg1:types.Connections):Promise<types.JSResp>;

export function Start(arg1:context.Context):Promise<void>;

export function TestConnection(arg1:types.ConnectionConfig):Promise<types.JSResp>;
