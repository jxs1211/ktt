// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {ai} from '../models';
import {types} from '../models';
import {context} from '../models';

export function Configure(arg1:ai.AIProvider):Promise<types.JSResp>;

export function GetAIProviders():Promise<types.JSResp>;

export function GetCompletion(arg1:string,arg2:string):Promise<types.JSResp>;

export function GetCompletion2(arg1:string,arg2:string,arg3:string):Promise<void>;

export function SetContext(arg1:context.Context):Promise<void>;