// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {types} from '../models';
import {context} from '../models';

export function Analyze(arg1:string,arg2:string,arg3:string,arg4:string,arg5:Array<string>,arg6:boolean,arg7:boolean,arg8:boolean):Promise<types.JSResp>;

export function CurrentContext():Promise<string>;

export function GetAvailableFilteredResources():Promise<types.JSResp>;

export function GetClusters():Promise<Array<string>>;

export function GetContexts():Promise<types.JSResp>;

export function GetLocalConfig():Promise<types.JSResp>;

export function LoadConfig(arg1:string):Promise<types.JSResp>;

export function Start(arg1:context.Context):Promise<void>;

export function TestConnection(arg1:string):Promise<types.JSResp>;
