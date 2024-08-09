export namespace convutil {
	
	export class CmdConvert {
	
	
	    static createFrom(source: any = {}) {
	        return new CmdConvert(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

export namespace types {
	
	export class Backend {
	    name: string;
	    model: string;
	    baseUrl: string;
	    appId: string;
	    apiKey: string;
	
	    static createFrom(source: any = {}) {
	        return new Backend(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.model = source["model"];
	        this.baseUrl = source["baseUrl"];
	        this.appId = source["appId"];
	        this.apiKey = source["apiKey"];
	    }
	}
	export class ConnectionProxy {
	    type?: number;
	    schema?: string;
	    addr?: string;
	    port?: number;
	    username?: string;
	    password?: string;
	
	    static createFrom(source: any = {}) {
	        return new ConnectionProxy(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.schema = source["schema"];
	        this.addr = source["addr"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	    }
	}
	export class ConnectionCluster {
	    enable?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ConnectionCluster(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enable = source["enable"];
	    }
	}
	export class ConnectionSentinel {
	    enable?: boolean;
	    master?: string;
	    username?: string;
	    password?: string;
	
	    static createFrom(source: any = {}) {
	        return new ConnectionSentinel(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enable = source["enable"];
	        this.master = source["master"];
	        this.username = source["username"];
	        this.password = source["password"];
	    }
	}
	export class ConnectionSSH {
	    enable?: boolean;
	    addr?: string;
	    port?: number;
	    loginType?: string;
	    username?: string;
	    password?: string;
	    pkFile?: string;
	    passphrase?: string;
	
	    static createFrom(source: any = {}) {
	        return new ConnectionSSH(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enable = source["enable"];
	        this.addr = source["addr"];
	        this.port = source["port"];
	        this.loginType = source["loginType"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.pkFile = source["pkFile"];
	        this.passphrase = source["passphrase"];
	    }
	}
	export class ConnectionSSL {
	    enable?: boolean;
	    keyFile?: string;
	    certFile?: string;
	    caFile?: string;
	    allowInsecure?: boolean;
	    sni?: string;
	
	    static createFrom(source: any = {}) {
	        return new ConnectionSSL(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enable = source["enable"];
	        this.keyFile = source["keyFile"];
	        this.certFile = source["certFile"];
	        this.caFile = source["caFile"];
	        this.allowInsecure = source["allowInsecure"];
	        this.sni = source["sni"];
	    }
	}
	export class Connection {
	    name: string;
	    group?: string;
	    lastDB: number;
	    network?: string;
	    sock?: string;
	    addr?: string;
	    port?: number;
	    username?: string;
	    password?: string;
	    defaultFilter?: string;
	    keySeparator?: string;
	    connTimeout?: number;
	    execTimeout?: number;
	    dbFilterType: string;
	    dbFilterList: number[];
	    keyView?: number;
	    loadSize?: number;
	    markColor?: string;
	    refreshInterval?: number;
	    alias?: {[key: number]: string};
	    ssl?: ConnectionSSL;
	    ssh?: ConnectionSSH;
	    sentinel?: ConnectionSentinel;
	    cluster?: ConnectionCluster;
	    proxy?: ConnectionProxy;
	    type?: string;
	    connections?: Connection[];
	
	    static createFrom(source: any = {}) {
	        return new Connection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.group = source["group"];
	        this.lastDB = source["lastDB"];
	        this.network = source["network"];
	        this.sock = source["sock"];
	        this.addr = source["addr"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.defaultFilter = source["defaultFilter"];
	        this.keySeparator = source["keySeparator"];
	        this.connTimeout = source["connTimeout"];
	        this.execTimeout = source["execTimeout"];
	        this.dbFilterType = source["dbFilterType"];
	        this.dbFilterList = source["dbFilterList"];
	        this.keyView = source["keyView"];
	        this.loadSize = source["loadSize"];
	        this.markColor = source["markColor"];
	        this.refreshInterval = source["refreshInterval"];
	        this.alias = source["alias"];
	        this.ssl = this.convertValues(source["ssl"], ConnectionSSL);
	        this.ssh = this.convertValues(source["ssh"], ConnectionSSH);
	        this.sentinel = this.convertValues(source["sentinel"], ConnectionSentinel);
	        this.cluster = this.convertValues(source["cluster"], ConnectionCluster);
	        this.proxy = this.convertValues(source["proxy"], ConnectionProxy);
	        this.type = source["type"];
	        this.connections = this.convertValues(source["connections"], Connection);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class ConnectionConfig {
	    name: string;
	    group?: string;
	    lastDB: number;
	    network?: string;
	    sock?: string;
	    addr?: string;
	    port?: number;
	    username?: string;
	    password?: string;
	    defaultFilter?: string;
	    keySeparator?: string;
	    connTimeout?: number;
	    execTimeout?: number;
	    dbFilterType: string;
	    dbFilterList: number[];
	    keyView?: number;
	    loadSize?: number;
	    markColor?: string;
	    refreshInterval?: number;
	    alias?: {[key: number]: string};
	    ssl?: ConnectionSSL;
	    ssh?: ConnectionSSH;
	    sentinel?: ConnectionSentinel;
	    cluster?: ConnectionCluster;
	    proxy?: ConnectionProxy;
	
	    static createFrom(source: any = {}) {
	        return new ConnectionConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.group = source["group"];
	        this.lastDB = source["lastDB"];
	        this.network = source["network"];
	        this.sock = source["sock"];
	        this.addr = source["addr"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.defaultFilter = source["defaultFilter"];
	        this.keySeparator = source["keySeparator"];
	        this.connTimeout = source["connTimeout"];
	        this.execTimeout = source["execTimeout"];
	        this.dbFilterType = source["dbFilterType"];
	        this.dbFilterList = source["dbFilterList"];
	        this.keyView = source["keyView"];
	        this.loadSize = source["loadSize"];
	        this.markColor = source["markColor"];
	        this.refreshInterval = source["refreshInterval"];
	        this.alias = source["alias"];
	        this.ssl = this.convertValues(source["ssl"], ConnectionSSL);
	        this.ssh = this.convertValues(source["ssh"], ConnectionSSH);
	        this.sentinel = this.convertValues(source["sentinel"], ConnectionSentinel);
	        this.cluster = this.convertValues(source["cluster"], ConnectionCluster);
	        this.proxy = this.convertValues(source["proxy"], ConnectionProxy);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	export class GetHashParam {
	    server: string;
	    db: number;
	    key: any;
	    field?: string;
	    format?: string;
	    decode?: string;
	
	    static createFrom(source: any = {}) {
	        return new GetHashParam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.server = source["server"];
	        this.db = source["db"];
	        this.key = source["key"];
	        this.field = source["field"];
	        this.format = source["format"];
	        this.decode = source["decode"];
	    }
	}
	export class JSResp {
	    success: boolean;
	    msg: string;
	    data?: any;
	
	    static createFrom(source: any = {}) {
	        return new JSResp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.msg = source["msg"];
	        this.data = source["data"];
	    }
	}
	export class KeyDetailParam {
	    server: string;
	    db: number;
	    key: any;
	    format?: string;
	    decode?: string;
	    matchPattern?: string;
	    reset: boolean;
	    full: boolean;
	
	    static createFrom(source: any = {}) {
	        return new KeyDetailParam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.server = source["server"];
	        this.db = source["db"];
	        this.key = source["key"];
	        this.format = source["format"];
	        this.decode = source["decode"];
	        this.matchPattern = source["matchPattern"];
	        this.reset = source["reset"];
	        this.full = source["full"];
	    }
	}
	export class KeySummaryParam {
	    server: string;
	    db: number;
	    key: any;
	
	    static createFrom(source: any = {}) {
	        return new KeySummaryParam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.server = source["server"];
	        this.db = source["db"];
	        this.key = source["key"];
	    }
	}
	export class PreferencesDecoder {
	    name: string;
	    enable: boolean;
	    auto: boolean;
	    decodePath: string;
	    decodeArgs: string[];
	    encodePath: string;
	    encodeArgs: string[];
	
	    static createFrom(source: any = {}) {
	        return new PreferencesDecoder(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.enable = source["enable"];
	        this.auto = source["auto"];
	        this.decodePath = source["decodePath"];
	        this.decodeArgs = source["decodeArgs"];
	        this.encodePath = source["encodePath"];
	        this.encodeArgs = source["encodeArgs"];
	    }
	}
	export class PreferencesCli {
	    fontFamily: string[];
	    fontSize: number;
	    cursorStyle: string;
	
	    static createFrom(source: any = {}) {
	        return new PreferencesCli(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fontFamily = source["fontFamily"];
	        this.fontSize = source["fontSize"];
	        this.cursorStyle = source["cursorStyle"];
	    }
	}
	export class PreferencesEditor {
	    font: string;
	    fontFamily: string[];
	    fontSize: number;
	    showLineNum: boolean;
	    showFolding: boolean;
	    dropText: boolean;
	    links: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PreferencesEditor(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.font = source["font"];
	        this.fontFamily = source["fontFamily"];
	        this.fontSize = source["fontSize"];
	        this.showLineNum = source["showLineNum"];
	        this.showFolding = source["showFolding"];
	        this.dropText = source["dropText"];
	        this.links = source["links"];
	    }
	}
	export class PreferencesAI {
	    enable: boolean;
	    explain: boolean;
	    backend: string;
	    backends: Backend[];
	
	    static createFrom(source: any = {}) {
	        return new PreferencesAI(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enable = source["enable"];
	        this.explain = source["explain"];
	        this.backend = source["backend"];
	        this.backends = this.convertValues(source["backends"], Backend);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PreferencesGeneral {
	    theme: string;
	    language: string;
	    font: string;
	    fontFamily: string[];
	    fontSize: number;
	    scanSize: number;
	    keyIconStyle: number;
	    useSysProxy: boolean;
	    useSysProxyHttp: boolean;
	    checkUpdate: boolean;
	    skipVersion: string;
	    allowTrack: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PreferencesGeneral(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.theme = source["theme"];
	        this.language = source["language"];
	        this.font = source["font"];
	        this.fontFamily = source["fontFamily"];
	        this.fontSize = source["fontSize"];
	        this.scanSize = source["scanSize"];
	        this.keyIconStyle = source["keyIconStyle"];
	        this.useSysProxy = source["useSysProxy"];
	        this.useSysProxyHttp = source["useSysProxyHttp"];
	        this.checkUpdate = source["checkUpdate"];
	        this.skipVersion = source["skipVersion"];
	        this.allowTrack = source["allowTrack"];
	    }
	}
	export class PreferencesBehavior {
	    welcomed: boolean;
	    asideWidth: number;
	    windowWidth: number;
	    windowHeight: number;
	    windowMaximised: boolean;
	    windowPosX: number;
	    windowPosY: number;
	
	    static createFrom(source: any = {}) {
	        return new PreferencesBehavior(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.welcomed = source["welcomed"];
	        this.asideWidth = source["asideWidth"];
	        this.windowWidth = source["windowWidth"];
	        this.windowHeight = source["windowHeight"];
	        this.windowMaximised = source["windowMaximised"];
	        this.windowPosX = source["windowPosX"];
	        this.windowPosY = source["windowPosY"];
	    }
	}
	export class Preferences {
	    behavior: PreferencesBehavior;
	    general: PreferencesGeneral;
	    ai: PreferencesAI;
	    editor: PreferencesEditor;
	    cli: PreferencesCli;
	    decoder: PreferencesDecoder[];
	
	    static createFrom(source: any = {}) {
	        return new Preferences(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.behavior = this.convertValues(source["behavior"], PreferencesBehavior);
	        this.general = this.convertValues(source["general"], PreferencesGeneral);
	        this.ai = this.convertValues(source["ai"], PreferencesAI);
	        this.editor = this.convertValues(source["editor"], PreferencesEditor);
	        this.cli = this.convertValues(source["cli"], PreferencesCli);
	        this.decoder = this.convertValues(source["decoder"], PreferencesDecoder);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	
	export class SetHashParam {
	    server: string;
	    db: number;
	    key: any;
	    field?: string;
	    newField?: string;
	    value: any;
	    format?: string;
	    decode?: string;
	    retFormat?: string;
	    retDecode?: string;
	
	    static createFrom(source: any = {}) {
	        return new SetHashParam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.server = source["server"];
	        this.db = source["db"];
	        this.key = source["key"];
	        this.field = source["field"];
	        this.newField = source["newField"];
	        this.value = source["value"];
	        this.format = source["format"];
	        this.decode = source["decode"];
	        this.retFormat = source["retFormat"];
	        this.retDecode = source["retDecode"];
	    }
	}
	export class SetKeyParam {
	    server: string;
	    db: number;
	    key: any;
	    keyType: string;
	    value: any;
	    ttl: number;
	    format?: string;
	    decode?: string;
	
	    static createFrom(source: any = {}) {
	        return new SetKeyParam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.server = source["server"];
	        this.db = source["db"];
	        this.key = source["key"];
	        this.keyType = source["keyType"];
	        this.value = source["value"];
	        this.ttl = source["ttl"];
	        this.format = source["format"];
	        this.decode = source["decode"];
	    }
	}
	export class SetListParam {
	    server: string;
	    db: number;
	    key: any;
	    index: number;
	    value: any;
	    format?: string;
	    decode?: string;
	    retFormat?: string;
	    retDecode?: string;
	
	    static createFrom(source: any = {}) {
	        return new SetListParam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.server = source["server"];
	        this.db = source["db"];
	        this.key = source["key"];
	        this.index = source["index"];
	        this.value = source["value"];
	        this.format = source["format"];
	        this.decode = source["decode"];
	        this.retFormat = source["retFormat"];
	        this.retDecode = source["retDecode"];
	    }
	}
	export class SetSetParam {
	    server: string;
	    db: number;
	    key: any;
	    value: any;
	    newValue: any;
	    format?: string;
	    decode?: string;
	    retFormat?: string;
	    retDecode?: string;
	
	    static createFrom(source: any = {}) {
	        return new SetSetParam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.server = source["server"];
	        this.db = source["db"];
	        this.key = source["key"];
	        this.value = source["value"];
	        this.newValue = source["newValue"];
	        this.format = source["format"];
	        this.decode = source["decode"];
	        this.retFormat = source["retFormat"];
	        this.retDecode = source["retDecode"];
	    }
	}
	export class SetZSetParam {
	    server: string;
	    db: number;
	    key: any;
	    value: any;
	    newValue: any;
	    score: number;
	    format?: string;
	    decode?: string;
	    retFormat?: string;
	    retDecode?: string;
	
	    static createFrom(source: any = {}) {
	        return new SetZSetParam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.server = source["server"];
	        this.db = source["db"];
	        this.key = source["key"];
	        this.value = source["value"];
	        this.newValue = source["newValue"];
	        this.score = source["score"];
	        this.format = source["format"];
	        this.decode = source["decode"];
	        this.retFormat = source["retFormat"];
	        this.retDecode = source["retDecode"];
	    }
	}

}

