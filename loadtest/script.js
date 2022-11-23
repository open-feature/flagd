import http from 'k6/http';
import {sleep} from 'k6';

export const options = {
    stages: [
        { duration: '1m', target: 200 },
        { duration: '3m', target: 200},
        { duration: '30s', target: 0 },
    ]
}

export const customOptions = {
    ffCount: 2500,
    type: "boolean"
}

export const prefix = "flag"

export default function () {
    let flag = prefix + Math.floor((Math.random() * customOptions.ffCount))

    let resp = http.post(genUrl(customOptions.type),
        JSON.stringify({
            flagKey: flag,
            context: {}
        }),
        {headers: {'Content-Type': 'application/json'}});

    console.log(JSON.stringify(resp.body))

    sleep(2);
}

export function genUrl(type) {
    switch (type) {
        case "boolean":
            return "http://localhost:8013/schema.v1.Service/ResolveBoolean"
        case "string":
            return "http://localhost:8013/schema.v1.Service/ResolveString"
    }

}