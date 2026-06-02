export namespace models {
	
	export class Collection {
	    id: number;
	    project_id: number;
	    parent_id?: number;
	    name: string;
	    sort_order: number;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Collection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.project_id = source["project_id"];
	        this.parent_id = source["parent_id"];
	        this.name = source["name"];
	        this.sort_order = source["sort_order"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}
	export class Cookie {
	    name: string;
	    value: string;
	    domain: string;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new Cookie(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.value = source["value"];
	        this.domain = source["domain"];
	        this.path = source["path"];
	    }
	}
	export class EnvVariable {
	    id: number;
	    environment_id: number;
	    key: string;
	    value: string;
	    enabled: boolean;
	    sort_order: number;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new EnvVariable(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.environment_id = source["environment_id"];
	        this.key = source["key"];
	        this.value = source["value"];
	        this.enabled = source["enabled"];
	        this.sort_order = source["sort_order"];
	        this.created_at = source["created_at"];
	    }
	}
	export class Environment {
	    id: number;
	    project_id: number;
	    name: string;
	    base_url: string;
	    is_active: boolean;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Environment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.project_id = source["project_id"];
	        this.name = source["name"];
	        this.base_url = source["base_url"];
	        this.is_active = source["is_active"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}
	export class HTTPResponse {
	    status: number;
	    status_text: string;
	    time: number;
	    size: number;
	    headers: Record<string, string>;
	    cookies: Cookie[];
	    body: string;
	    raw_request: string;
	    curl_command: string;
	
	    static createFrom(source: any = {}) {
	        return new HTTPResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.status_text = source["status_text"];
	        this.time = source["time"];
	        this.size = source["size"];
	        this.headers = source["headers"];
	        this.cookies = this.convertValues(source["cookies"], Cookie);
	        this.body = source["body"];
	        this.raw_request = source["raw_request"];
	        this.curl_command = source["curl_command"];
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
	export class History {
	    id: number;
	    project_id: number;
	    request_id?: number;
	    method: string;
	    url: string;
	    request_headers: string;
	    request_body: string;
	    response_status: number;
	    response_headers: string;
	    response_body: string;
	    duration_ms: number;
	    created_at: string;
	
	    static createFrom(source: any = {}) {
	        return new History(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.project_id = source["project_id"];
	        this.request_id = source["request_id"];
	        this.method = source["method"];
	        this.url = source["url"];
	        this.request_headers = source["request_headers"];
	        this.request_body = source["request_body"];
	        this.response_status = source["response_status"];
	        this.response_headers = source["response_headers"];
	        this.response_body = source["response_body"];
	        this.duration_ms = source["duration_ms"];
	        this.created_at = source["created_at"];
	    }
	}
	export class Project {
	    id: number;
	    name: string;
	    description: string;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}
	export class ProjectStats {
	    request_count: number;
	    collection_count: number;
	
	    static createFrom(source: any = {}) {
	        return new ProjectStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.request_count = source["request_count"];
	        this.collection_count = source["collection_count"];
	    }
	}
	export class Request {
	    id: number;
	    collection_id: number;
	    name: string;
	    description: string;
	    method: string;
	    url: string;
	    headers: string;
	    params: string;
	    body_type: string;
	    body: string;
	    auth: string;
	    sort_order: number;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new Request(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.collection_id = source["collection_id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.method = source["method"];
	        this.url = source["url"];
	        this.headers = source["headers"];
	        this.params = source["params"];
	        this.body_type = source["body_type"];
	        this.body = source["body"];
	        this.auth = source["auth"];
	        this.sort_order = source["sort_order"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}
	export class TreeItem {
	    id: number;
	    name: string;
	    type: string;
	    method?: string;
	    url?: string;
	    children: TreeItem[];
	    sort_order: number;
	
	    static createFrom(source: any = {}) {
	        return new TreeItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.method = source["method"];
	        this.url = source["url"];
	        this.children = this.convertValues(source["children"], TreeItem);
	        this.sort_order = source["sort_order"];
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

export namespace services {
	
	export class ImportResult {
	    collections: number;
	    requests: number;
	
	    static createFrom(source: any = {}) {
	        return new ImportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.collections = source["collections"];
	        this.requests = source["requests"];
	    }
	}

}

