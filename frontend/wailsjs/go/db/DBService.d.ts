// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {types} from '../models';
import {context} from '../models';

export function CreateSession(arg1:string,arg2:string,arg3:string,arg4:string):Promise<types.JSResp>;

export function DeleteSession(arg1:number):Promise<types.JSResp>;

export function GetSessionsByClusterName(arg1:string):Promise<types.JSResp>;

export function Start(arg1:context.Context):Promise<void>;
