import http from 'k6/http';

/*
* Sample K6 (https://k6.io/) load test script
* */

// K6 options - Load generation pattern: Ramp up, hold and teardown
export const options = {
    stages: [{duration: '10s', target: 50}, {duration: '30s', target: 50}, {duration: '10s', target: 0},]
}

// Flag prefix - See ff_gen.go to match
export const prefix = "flag"

// Custom options : Number of FFs flagd serves and type of the FFs being served
export const customOptions = {
    ffCount: 100,
    type: "boolean"
}

export default function () {
    // Randomly select flag to evaluate
    let flag = prefix + Math.floor((Math.random() * customOptions.ffCount))

    let resp = http.post(genUrl(customOptions.type), JSON.stringify({
        flagKey: flag, context: {}
    }), {headers: {'Content-Type': 'application/json'}});

    // Handle and report errors
    if (resp.status !== 200) {
        console.log("Error response - FlagId : " + flag + " Response :" + JSON.stringify(resp.body))
    }
}

export function genUrl(type) {
    switch (type) {
        case "boolean":
            return "http://localhost:8013/schema.v1.Service/ResolveBoolean"
        case "string":
            return "http://localhost:8013/schema.v1.Service/ResolveString"
    }
}