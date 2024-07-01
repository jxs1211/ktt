export namespace types {
	
	export class Connection {
	    group: string;
	    name: string;
	    addr: string;
	    port: number;
	    username: string;
	    password: string;
	    defaultFilter: string;
	    keySeparator: string;
	    connTimeout: number;
	    execTimeout: number;
	    markColor: string;
	
	    static createFrom(source: any = {}) {
	        return new Connection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.group = source["group"];
	        this.name = source["name"];
	        this.addr = source["addr"];
	        this.port = source["port"];
	        this.username = source["username"];
	        this.password = source["password"];
	        this.defaultFilter = source["defaultFilter"];
	        this.keySeparator = source["keySeparator"];
	        this.connTimeout = source["connTimeout"];
	        this.execTimeout = source["execTimeout"];
	        this.markColor = source["markColor"];
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

}

