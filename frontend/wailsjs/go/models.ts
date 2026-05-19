export namespace email {
	
	export class DirectMailConfig {
	    name: string;
	    baseUrl: string;
	    refreshToken: string;
	    clientId: string;
	    email: string;
	    mailbox: string;
	
	    static createFrom(source: any = {}) {
	        return new DirectMailConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.baseUrl = source["baseUrl"];
	        this.refreshToken = source["refreshToken"];
	        this.clientId = source["clientId"];
	        this.email = source["email"];
	        this.mailbox = source["mailbox"];
	    }
	}
	export class DuckDuckGoConfig {
	    token: string;
	
	    static createFrom(source: any = {}) {
	        return new DuckDuckGoConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.token = source["token"];
	    }
	}
	export class MoeMailConfig {
	    name: string;
	    url: string;
	    apiKey: string;
	
	    static createFrom(source: any = {}) {
	        return new MoeMailConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.url = source["url"];
	        this.apiKey = source["apiKey"];
	    }
	}
	export class TEmailConfig {
	    name: string;
	    baseUrl: string;
	    email: string;
	    jwt?: string;
	    adminPassword?: string;
	
	    static createFrom(source: any = {}) {
	        return new TEmailConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.baseUrl = source["baseUrl"];
	        this.email = source["email"];
	        this.jwt = source["jwt"];
	        this.adminPassword = source["adminPassword"];
	    }
	}

}

export namespace proxy {
	
	export class Info {
	    ok: boolean;
	    scheme: string;
	    ip: string;
	    country: string;
	    region: string;
	    city: string;
	    isp: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.scheme = source["scheme"];
	        this.ip = source["ip"];
	        this.country = source["country"];
	        this.region = source["region"];
	        this.city = source["city"];
	        this.isp = source["isp"];
	        this.error = source["error"];
	    }
	}

}

export namespace task {
	
	export class StartTaskRequest {
	    count: number;
	    concurrency: number;
	    delay: number;
	    outputPath: string;
	    emailProvider: string;
	    moemailDomains: string[];
	    moemailConfigs: Record<string, Array<email.MoeMailConfig>>;
	    moemailRandomMode: boolean;
	    duckTokens: string[];
	    temailConfigs: email.TEmailConfig[];
	    directMailConfigs: email.DirectMailConfig[];
	
	    static createFrom(source: any = {}) {
	        return new StartTaskRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.count = source["count"];
	        this.concurrency = source["concurrency"];
	        this.delay = source["delay"];
	        this.outputPath = source["outputPath"];
	        this.emailProvider = source["emailProvider"];
	        this.moemailDomains = source["moemailDomains"];
	        this.moemailConfigs = this.convertValues(source["moemailConfigs"], Array<email.MoeMailConfig>, true);
	        this.moemailRandomMode = source["moemailRandomMode"];
	        this.duckTokens = source["duckTokens"];
	        this.temailConfigs = this.convertValues(source["temailConfigs"], email.TEmailConfig);
	        this.directMailConfigs = this.convertValues(source["directMailConfigs"], email.DirectMailConfig);
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

}

