export namespace handlers {
	
	export class AssertRule {
	    type: string;
	    target: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new AssertRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.target = source["target"];
	        this.value = source["value"];
	    }
	}
	export class AssertResult {
	    rule: AssertRule;
	    passed: boolean;
	    actual: string;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new AssertResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.rule = this.convertValues(source["rule"], AssertRule);
	        this.passed = source["passed"];
	        this.actual = source["actual"];
	        this.error = source["error"];
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
	
	export class CodegenInput {
	    Method: string;
	    URL: string;
	    Headers: Record<string, string>;
	    Body: string;
	    BodyType: string;
	
	    static createFrom(source: any = {}) {
	        return new CodegenInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Method = source["Method"];
	        this.URL = source["URL"];
	        this.Headers = source["Headers"];
	        this.Body = source["Body"];
	        this.BodyType = source["BodyType"];
	    }
	}
	export class ImportCollection {
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new ImportCollection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	    }
	}
	export class ImportRequest {
	    name: string;
	    collection_id: string;
	    method: string;
	    url: string;
	    headers: string;
	    params: string;
	    body: string;
	    auth: string;
	
	    static createFrom(source: any = {}) {
	        return new ImportRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.collection_id = source["collection_id"];
	        this.method = source["method"];
	        this.url = source["url"];
	        this.headers = source["headers"];
	        this.params = source["params"];
	        this.body = source["body"];
	        this.auth = source["auth"];
	    }
	}
	export class ImportResponse {
	    collections: ImportCollection[];
	    requests: ImportRequest[];
	
	    static createFrom(source: any = {}) {
	        return new ImportResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.collections = this.convertValues(source["collections"], ImportCollection);
	        this.requests = this.convertValues(source["requests"], ImportRequest);
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
	export class RecordHistoryInput {
	    ProjectID: string;
	    RequestID: string;
	    Method: string;
	    URL: string;
	    Headers: string;
	    Body: string;
	    ResponseStatus: number;
	    ResponseBody: string;
	    ResponseHeaders: string;
	    DurationMs: number;
	
	    static createFrom(source: any = {}) {
	        return new RecordHistoryInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ProjectID = source["ProjectID"];
	        this.RequestID = source["RequestID"];
	        this.Method = source["Method"];
	        this.URL = source["URL"];
	        this.Headers = source["Headers"];
	        this.Body = source["Body"];
	        this.ResponseStatus = source["ResponseStatus"];
	        this.ResponseBody = source["ResponseBody"];
	        this.ResponseHeaders = source["ResponseHeaders"];
	        this.DurationMs = source["DurationMs"];
	    }
	}
	export class RunAssertsInput {
	    Method: string;
	    URL: string;
	    Headers: Record<string, string>;
	    Body: string;
	    BodyType: string;
	    BodyFiles: httpclient.BodyFile[];
	    AuthType: string;
	    AuthData: Record<string, string>;
	    TimeoutMs: number;
	    FollowRedirect: boolean;
	    Asserts: AssertRule[];
	
	    static createFrom(source: any = {}) {
	        return new RunAssertsInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Method = source["Method"];
	        this.URL = source["URL"];
	        this.Headers = source["Headers"];
	        this.Body = source["Body"];
	        this.BodyType = source["BodyType"];
	        this.BodyFiles = this.convertValues(source["BodyFiles"], httpclient.BodyFile);
	        this.AuthType = source["AuthType"];
	        this.AuthData = source["AuthData"];
	        this.TimeoutMs = source["TimeoutMs"];
	        this.FollowRedirect = source["FollowRedirect"];
	        this.Asserts = this.convertValues(source["Asserts"], AssertRule);
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
	export class RunAssertsResponse {
	    response?: httpclient.Response;
	    asserts: AssertResult[];
	
	    static createFrom(source: any = {}) {
	        return new RunAssertsResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.response = this.convertValues(source["response"], httpclient.Response);
	        this.asserts = this.convertValues(source["asserts"], AssertResult);
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
	export class SendRequestInput {
	    Method: string;
	    URL: string;
	    Headers: Record<string, string>;
	    Body: string;
	    BodyType: string;
	    BodyFiles: httpclient.BodyFile[];
	    AuthType: string;
	    AuthData: Record<string, string>;
	    TimeoutMs: number;
	    FollowRedirect: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SendRequestInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Method = source["Method"];
	        this.URL = source["URL"];
	        this.Headers = source["Headers"];
	        this.Body = source["Body"];
	        this.BodyType = source["BodyType"];
	        this.BodyFiles = this.convertValues(source["BodyFiles"], httpclient.BodyFile);
	        this.AuthType = source["AuthType"];
	        this.AuthData = source["AuthData"];
	        this.TimeoutMs = source["TimeoutMs"];
	        this.FollowRedirect = source["FollowRedirect"];
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

export namespace httpclient {
	
	export class BodyFile {
	    key: string;
	    value: string;
	    file_path: string;
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new BodyFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.value = source["value"];
	        this.file_path = source["file_path"];
	        this.enabled = source["enabled"];
	    }
	}
	export class Response {
	    status: number;
	    status_text: string;
	    headers: Record<string, Array<string>>;
	    body: string;
	    duration_ms: number;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.status_text = source["status_text"];
	        this.headers = source["headers"];
	        this.body = source["body"];
	        this.duration_ms = source["duration_ms"];
	    }
	}

}

export namespace models {
	
	export class Collection {
	    id: string;
	    project_id: string;
	    parent_id?: string;
	    name: string;
	    sort_order: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
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
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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
	export class Environment {
	    id: string;
	    project_id: string;
	    name: string;
	    variables: string;
	    is_active: boolean;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Environment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.project_id = source["project_id"];
	        this.name = source["name"];
	        this.variables = source["variables"];
	        this.is_active = source["is_active"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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
	    id: string;
	    project_id: string;
	    request_id?: string;
	    method: string;
	    url: string;
	    headers: string;
	    body: string;
	    response_status: number;
	    response_body: string;
	    response_headers: string;
	    duration_ms: number;
	    // Go type: time
	    created_at: any;
	
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
	        this.headers = source["headers"];
	        this.body = source["body"];
	        this.response_status = source["response_status"];
	        this.response_body = source["response_body"];
	        this.response_headers = source["response_headers"];
	        this.duration_ms = source["duration_ms"];
	        this.created_at = this.convertValues(source["created_at"], null);
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
	export class Project {
	    id: string;
	    name: string;
	    description: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Project(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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
	export class Request {
	    id: string;
	    collection_id: string;
	    name: string;
	    method: string;
	    url: string;
	    headers: string;
	    params: string;
	    body: string;
	    auth: string;
	    script: string;
	    sort_order: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Request(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.collection_id = source["collection_id"];
	        this.name = source["name"];
	        this.method = source["method"];
	        this.url = source["url"];
	        this.headers = source["headers"];
	        this.params = source["params"];
	        this.body = source["body"];
	        this.auth = source["auth"];
	        this.script = source["script"];
	        this.sort_order = source["sort_order"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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
	
	export class CookieInfo {
	    domain: string;
	    name: string;
	    value: string;
	    path: string;
	    // Go type: time
	    expires: any;
	    secure: boolean;
	    http_only: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CookieInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.domain = source["domain"];
	        this.name = source["name"];
	        this.value = source["value"];
	        this.path = source["path"];
	        this.expires = this.convertValues(source["expires"], null);
	        this.secure = source["secure"];
	        this.http_only = source["http_only"];
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

