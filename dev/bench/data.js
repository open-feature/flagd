window.BENCHMARK_DATA = {
  "lastUpdate": 1674613672230,
  "repoUrl": "https://github.com/open-feature/flagd",
  "entries": {
    "Go Benchmark": [
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "5201f6b753c7e663c29343e76df96767511c78e6",
          "message": "chore: fixing multi-arch build failure (#153)\n\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>",
          "timestamp": "2022-08-22T18:03:24Z",
          "url": "https://github.com/open-feature/flagd/commit/5201f6b753c7e663c29343e76df96767511c78e6"
        },
        "date": 1661223171368,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2217,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "517521 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 6460,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "178453 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2214,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "522806 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 6546,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "181520 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 6488,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "179304 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2218,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "521745 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2229,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "518455 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 6525,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "178917 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3592,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "301960 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 6275,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "184858 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c72323eb5a3f371a64301fe4c10a368705d6e8e9",
          "message": "ci: Benchmark threshold update (#154)\n\nSigned-off-by: James-Milligan <james@omnant.co.uk>",
          "timestamp": "2022-08-24T00:58:01Z",
          "url": "https://github.com/open-feature/flagd/commit/c72323eb5a3f371a64301fe4c10a368705d6e8e9"
        },
        "date": 1661309425536,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2144,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "530707 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 6271,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "183297 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2160,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "523810 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 6270,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "180550 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2161,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "545985 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 6268,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "184628 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2151,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "541696 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 6354,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "187074 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3464,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "341629 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 6223,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "172928 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c72323eb5a3f371a64301fe4c10a368705d6e8e9",
          "message": "ci: Benchmark threshold update (#154)\n\nSigned-off-by: James-Milligan <james@omnant.co.uk>",
          "timestamp": "2022-08-24T00:58:01Z",
          "url": "https://github.com/open-feature/flagd/commit/c72323eb5a3f371a64301fe4c10a368705d6e8e9"
        },
        "date": 1661395918957,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2145,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "536017 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 6257,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "181125 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2144,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "541554 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 6227,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "188649 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2164,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "548763 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 6253,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "187670 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2143,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "547308 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 6259,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "189404 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3726,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "342313 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 6144,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "185961 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "17ef4c69bbe1b4b4acb43481d0461ecead57bbb2",
          "message": "fix: bug in the logger setup (#156)\n\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>\r\n\r\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>",
          "timestamp": "2022-08-25T16:28:13Z",
          "url": "https://github.com/open-feature/flagd/commit/17ef4c69bbe1b4b4acb43481d0461ecead57bbb2"
        },
        "date": 1661482480579,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 1616,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "718870 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 4792,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "244371 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 1619,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "730627 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 4868,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "238924 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 1623,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "735466 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 5615,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "239529 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 1673,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "600645 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 4874,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "240303 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 2614,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "416019 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 4668,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "248421 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "17ef4c69bbe1b4b4acb43481d0461ecead57bbb2",
          "message": "fix: bug in the logger setup (#156)\n\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>\r\n\r\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>",
          "timestamp": "2022-08-25T16:28:13Z",
          "url": "https://github.com/open-feature/flagd/commit/17ef4c69bbe1b4b4acb43481d0461ecead57bbb2"
        },
        "date": 1661568681467,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2603,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "433876 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 8470,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "133934 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2725,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "437341 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 8908,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "142506 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2935,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "383500 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 8658,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "140103 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2838,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "378596 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 8492,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "152095 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 4465,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "284115 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 9099,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "145290 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "17ef4c69bbe1b4b4acb43481d0461ecead57bbb2",
          "message": "fix: bug in the logger setup (#156)\n\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>\r\n\r\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>",
          "timestamp": "2022-08-25T16:28:13Z",
          "url": "https://github.com/open-feature/flagd/commit/17ef4c69bbe1b4b4acb43481d0461ecead57bbb2"
        },
        "date": 1661655184587,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2653,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "431786 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 7549,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "160660 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2777,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "447568 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 7860,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "143872 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2676,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "454460 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 7642,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "142744 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2673,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "451774 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 7788,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "132781 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 4248,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "265743 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 7468,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "157278 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "17ef4c69bbe1b4b4acb43481d0461ecead57bbb2",
          "message": "fix: bug in the logger setup (#156)\n\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>\r\n\r\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>",
          "timestamp": "2022-08-25T16:28:13Z",
          "url": "https://github.com/open-feature/flagd/commit/17ef4c69bbe1b4b4acb43481d0461ecead57bbb2"
        },
        "date": 1661741763378,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 6520,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "181639 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2205,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "540945 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2201,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "512493 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 6374,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "182295 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2214,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "532401 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 6432,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "178321 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2186,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "539389 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 6711,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "185469 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3550,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "334668 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 6259,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "192130 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "17ef4c69bbe1b4b4acb43481d0461ecead57bbb2",
          "message": "fix: bug in the logger setup (#156)\n\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>\r\n\r\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>",
          "timestamp": "2022-08-25T16:28:13Z",
          "url": "https://github.com/open-feature/flagd/commit/17ef4c69bbe1b4b4acb43481d0461ecead57bbb2"
        },
        "date": 1661828093517,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2169,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "541860 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 6365,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "182014 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2183,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "532588 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 6429,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "185800 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2193,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "536719 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 6354,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "185162 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2174,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "540945 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 6401,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "186177 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3701,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "325796 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 6157,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "192830 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "17ef4c69bbe1b4b4acb43481d0461ecead57bbb2",
          "message": "fix: bug in the logger setup (#156)\n\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>\r\n\r\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>",
          "timestamp": "2022-08-25T16:28:13Z",
          "url": "https://github.com/open-feature/flagd/commit/17ef4c69bbe1b4b4acb43481d0461ecead57bbb2"
        },
        "date": 1661914639609,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2246,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "512008 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 6388,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "181123 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 6380,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "183691 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2169,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "537919 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2169,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "536521 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 6302,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "180568 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2154,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "545298 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 6308,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "184984 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3508,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "288982 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 6638,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "177416 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "17ef4c69bbe1b4b4acb43481d0461ecead57bbb2",
          "message": "fix: bug in the logger setup (#156)\n\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>\r\n\r\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>",
          "timestamp": "2022-08-25T16:28:13Z",
          "url": "https://github.com/open-feature/flagd/commit/17ef4c69bbe1b4b4acb43481d0461ecead57bbb2"
        },
        "date": 1662000726486,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2208,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "519520 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 6436,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "179425 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2225,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "520490 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 6503,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "178956 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2218,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "523653 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 6597,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "181119 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 6706,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "159484 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2222,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "523438 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3532,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "341924 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 6245,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "186378 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c977eb74020aafada587473db72c12356abc373c",
          "message": "ci: Increased benchmark test duration from 1s to 5s (#158)\n\nSigned-off-by: James-Milligan <james@omnant.co.uk>\r\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-01T12:48:10Z",
          "url": "https://github.com/open-feature/flagd/commit/c977eb74020aafada587473db72c12356abc373c"
        },
        "date": 1662087344655,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2639,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2289062 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 7702,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "759217 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2627,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2239776 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 7663,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "787921 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2617,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2313598 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 7891,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "736087 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2619,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2280148 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 7737,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "750913 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 4172,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1463862 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 7442,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "786360 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c977eb74020aafada587473db72c12356abc373c",
          "message": "ci: Increased benchmark test duration from 1s to 5s (#158)\n\nSigned-off-by: James-Milligan <james@omnant.co.uk>\r\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-01T12:48:10Z",
          "url": "https://github.com/open-feature/flagd/commit/c977eb74020aafada587473db72c12356abc373c"
        },
        "date": 1662173625355,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 7911,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "742731 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2656,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2206438 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2716,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2249876 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 7845,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "758742 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2660,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2272032 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 7788,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "709282 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2665,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2275237 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 7850,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "785713 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 4170,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1444104 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 7298,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "788148 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c977eb74020aafada587473db72c12356abc373c",
          "message": "ci: Increased benchmark test duration from 1s to 5s (#158)\n\nSigned-off-by: James-Milligan <james@omnant.co.uk>\r\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-01T12:48:10Z",
          "url": "https://github.com/open-feature/flagd/commit/c977eb74020aafada587473db72c12356abc373c"
        },
        "date": 1662260059916,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 3019,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2075259 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 9548,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "628844 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2909,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2006661 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 9421,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "669991 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2900,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2067385 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 9325,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "603486 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2898,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2047044 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 9298,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "667255 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 4528,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1341201 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 9179,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "643533 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c977eb74020aafada587473db72c12356abc373c",
          "message": "ci: Increased benchmark test duration from 1s to 5s (#158)\n\nSigned-off-by: James-Milligan <james@omnant.co.uk>\r\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-01T12:48:10Z",
          "url": "https://github.com/open-feature/flagd/commit/c977eb74020aafada587473db72c12356abc373c"
        },
        "date": 1662346504473,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2632,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2293872 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 7538,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "718416 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2679,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2309421 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 7572,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "752828 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2614,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2290938 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 7615,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "775242 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2562,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2298822 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 7539,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "761265 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 4008,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1516622 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 7170,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "817736 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c977eb74020aafada587473db72c12356abc373c",
          "message": "ci: Increased benchmark test duration from 1s to 5s (#158)\n\nSigned-off-by: James-Milligan <james@omnant.co.uk>\r\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-01T12:48:10Z",
          "url": "https://github.com/open-feature/flagd/commit/c977eb74020aafada587473db72c12356abc373c"
        },
        "date": 1662432986422,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2574,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2337337 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 6625,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "811669 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2382,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2518514 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 6772,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "811916 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 6736,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "949441 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2518,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2655914 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2493,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2531041 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 7442,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "736524 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3828,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1555594 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 6959,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "807177 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c4e030ed651ab800aa56a555471da5913eb95259",
          "message": "feat: kubernetes sync (#157)\n\nhttps://github.com/open-feature/research/pull/31\r\n<img width=\"1512\" alt=\"Screenshot 2022-08-30 at 15 54 16\" src=\"https://user-images.githubusercontent.com/1235925/187469966-d3c137bd-1f1b-4ee7-9512-848951ec780b.png\">\r\n\r\nSigned-off-by: James-Milligan <james@omnant.co.uk>\r\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>\r\nSigned-off-by: Alex Jones <alex.jones@canonical.com>\r\nCo-authored-by: James Milligan <75740990+james-milligan@users.noreply.github.com>\r\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-06T12:45:19Z",
          "url": "https://github.com/open-feature/flagd/commit/c4e030ed651ab800aa56a555471da5913eb95259"
        },
        "date": 1662519534049,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2581,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2302611 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 7656,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "775750 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 7609,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "757918 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2659,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2272206 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2644,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2262958 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 7620,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "770940 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2603,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2294538 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 7639,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "759272 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 7328,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "784206 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 4121,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1466067 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c4e030ed651ab800aa56a555471da5913eb95259",
          "message": "feat: kubernetes sync (#157)\n\nhttps://github.com/open-feature/research/pull/31\r\n<img width=\"1512\" alt=\"Screenshot 2022-08-30 at 15 54 16\" src=\"https://user-images.githubusercontent.com/1235925/187469966-d3c137bd-1f1b-4ee7-9512-848951ec780b.png\">\r\n\r\nSigned-off-by: James-Milligan <james@omnant.co.uk>\r\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>\r\nSigned-off-by: Alex Jones <alex.jones@canonical.com>\r\nCo-authored-by: James Milligan <75740990+james-milligan@users.noreply.github.com>\r\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-06T12:45:19Z",
          "url": "https://github.com/open-feature/flagd/commit/c4e030ed651ab800aa56a555471da5913eb95259"
        },
        "date": 1662605798983,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 6225,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "933054 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2158,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2797638 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2157,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2789206 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 6262,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "933214 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2147,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2788762 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 6231,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "922950 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2135,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2806914 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 6253,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "902455 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3329,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1808128 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 6026,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "973653 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c4e030ed651ab800aa56a555471da5913eb95259",
          "message": "feat: kubernetes sync (#157)\n\nhttps://github.com/open-feature/research/pull/31\r\n<img width=\"1512\" alt=\"Screenshot 2022-08-30 at 15 54 16\" src=\"https://user-images.githubusercontent.com/1235925/187469966-d3c137bd-1f1b-4ee7-9512-848951ec780b.png\">\r\n\r\nSigned-off-by: James-Milligan <james@omnant.co.uk>\r\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>\r\nSigned-off-by: Alex Jones <alex.jones@canonical.com>\r\nCo-authored-by: James Milligan <75740990+james-milligan@users.noreply.github.com>\r\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-06T12:45:19Z",
          "url": "https://github.com/open-feature/flagd/commit/c4e030ed651ab800aa56a555471da5913eb95259"
        },
        "date": 1662692201083,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2655,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2256415 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 7540,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "756637 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2546,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2357001 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 7653,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "789775 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2692,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2217356 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 7690,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "730690 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2649,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2325058 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 7562,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "777019 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 4065,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1467740 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 7084,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "807318 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c4e030ed651ab800aa56a555471da5913eb95259",
          "message": "feat: kubernetes sync (#157)\n\nhttps://github.com/open-feature/research/pull/31\r\n<img width=\"1512\" alt=\"Screenshot 2022-08-30 at 15 54 16\" src=\"https://user-images.githubusercontent.com/1235925/187469966-d3c137bd-1f1b-4ee7-9512-848951ec780b.png\">\r\n\r\nSigned-off-by: James-Milligan <james@omnant.co.uk>\r\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>\r\nSigned-off-by: Alex Jones <alex.jones@canonical.com>\r\nCo-authored-by: James Milligan <75740990+james-milligan@users.noreply.github.com>\r\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-06T12:45:19Z",
          "url": "https://github.com/open-feature/flagd/commit/c4e030ed651ab800aa56a555471da5913eb95259"
        },
        "date": 1662778597295,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2583,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2322351 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 7312,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "715806 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2552,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2382561 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 7480,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "782936 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2597,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2310476 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 7353,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "742916 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2567,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2328908 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 7394,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "732582 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 4031,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1519087 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 7318,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "810474 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c4e030ed651ab800aa56a555471da5913eb95259",
          "message": "feat: kubernetes sync (#157)\n\nhttps://github.com/open-feature/research/pull/31\r\n<img width=\"1512\" alt=\"Screenshot 2022-08-30 at 15 54 16\" src=\"https://user-images.githubusercontent.com/1235925/187469966-d3c137bd-1f1b-4ee7-9512-848951ec780b.png\">\r\n\r\nSigned-off-by: James-Milligan <james@omnant.co.uk>\r\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>\r\nSigned-off-by: Alex Jones <alex.jones@canonical.com>\r\nCo-authored-by: James Milligan <75740990+james-milligan@users.noreply.github.com>\r\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-06T12:45:19Z",
          "url": "https://github.com/open-feature/flagd/commit/c4e030ed651ab800aa56a555471da5913eb95259"
        },
        "date": 1662865002932,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2232,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2702359 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 6449,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "892836 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 6473,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "897904 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2230,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2700266 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2217,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2693624 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 6438,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "914001 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2212,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2711114 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 6416,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "910333 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 6239,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "940348 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3554,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1705554 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c4e030ed651ab800aa56a555471da5913eb95259",
          "message": "feat: kubernetes sync (#157)\n\nhttps://github.com/open-feature/research/pull/31\r\n<img width=\"1512\" alt=\"Screenshot 2022-08-30 at 15 54 16\" src=\"https://user-images.githubusercontent.com/1235925/187469966-d3c137bd-1f1b-4ee7-9512-848951ec780b.png\">\r\n\r\nSigned-off-by: James-Milligan <james@omnant.co.uk>\r\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>\r\nSigned-off-by: Alex Jones <alex.jones@canonical.com>\r\nCo-authored-by: James Milligan <75740990+james-milligan@users.noreply.github.com>\r\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-06T12:45:19Z",
          "url": "https://github.com/open-feature/flagd/commit/c4e030ed651ab800aa56a555471da5913eb95259"
        },
        "date": 1662951504260,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2474,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2437579 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 7287,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "733414 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 7752,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "801979 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2408,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2610741 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 7378,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "787598 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2372,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2442207 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 7363,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "855690 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2360,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2411420 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3453,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1644495 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 7392,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "843285 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c4e030ed651ab800aa56a555471da5913eb95259",
          "message": "feat: kubernetes sync (#157)\n\nhttps://github.com/open-feature/research/pull/31\r\n<img width=\"1512\" alt=\"Screenshot 2022-08-30 at 15 54 16\" src=\"https://user-images.githubusercontent.com/1235925/187469966-d3c137bd-1f1b-4ee7-9512-848951ec780b.png\">\r\n\r\nSigned-off-by: James-Milligan <james@omnant.co.uk>\r\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>\r\nSigned-off-by: Alex Jones <alex.jones@canonical.com>\r\nCo-authored-by: James Milligan <75740990+james-milligan@users.noreply.github.com>\r\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-06T12:45:19Z",
          "url": "https://github.com/open-feature/flagd/commit/c4e030ed651ab800aa56a555471da5913eb95259"
        },
        "date": 1663037860974,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 9199,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "622471 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2913,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2084194 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2879,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2081883 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 9131,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "603464 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2859,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2084920 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 9244,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "623262 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2833,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2051919 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 8904,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "659817 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 4363,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1364086 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 8782,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "667090 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c4e030ed651ab800aa56a555471da5913eb95259",
          "message": "feat: kubernetes sync (#157)\n\nhttps://github.com/open-feature/research/pull/31\r\n<img width=\"1512\" alt=\"Screenshot 2022-08-30 at 15 54 16\" src=\"https://user-images.githubusercontent.com/1235925/187469966-d3c137bd-1f1b-4ee7-9512-848951ec780b.png\">\r\n\r\nSigned-off-by: James-Milligan <james@omnant.co.uk>\r\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>\r\nSigned-off-by: Alex Jones <alex.jones@canonical.com>\r\nCo-authored-by: James Milligan <75740990+james-milligan@users.noreply.github.com>\r\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-06T12:45:19Z",
          "url": "https://github.com/open-feature/flagd/commit/c4e030ed651ab800aa56a555471da5913eb95259"
        },
        "date": 1663124167769,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2183,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2750144 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 6256,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "917758 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 6287,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "936871 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2183,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2755948 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2179,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2753133 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 6272,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "925438 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2156,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2775810 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 6287,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "934234 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3359,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1800811 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 6088,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "961833 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c4e030ed651ab800aa56a555471da5913eb95259",
          "message": "feat: kubernetes sync (#157)\n\nhttps://github.com/open-feature/research/pull/31\r\n<img width=\"1512\" alt=\"Screenshot 2022-08-30 at 15 54 16\" src=\"https://user-images.githubusercontent.com/1235925/187469966-d3c137bd-1f1b-4ee7-9512-848951ec780b.png\">\r\n\r\nSigned-off-by: James-Milligan <james@omnant.co.uk>\r\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>\r\nSigned-off-by: Alex Jones <alex.jones@canonical.com>\r\nCo-authored-by: James Milligan <75740990+james-milligan@users.noreply.github.com>\r\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-06T12:45:19Z",
          "url": "https://github.com/open-feature/flagd/commit/c4e030ed651ab800aa56a555471da5913eb95259"
        },
        "date": 1663210761918,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2998,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2039781 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 9367,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "628688 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2948,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2047084 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 9374,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "632007 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2927,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2012560 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 9352,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "582133 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2920,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "1995751 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 9303,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "639338 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 4431,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1355011 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 9136,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "638581 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c4e030ed651ab800aa56a555471da5913eb95259",
          "message": "feat: kubernetes sync (#157)\n\nhttps://github.com/open-feature/research/pull/31\r\n<img width=\"1512\" alt=\"Screenshot 2022-08-30 at 15 54 16\" src=\"https://user-images.githubusercontent.com/1235925/187469966-d3c137bd-1f1b-4ee7-9512-848951ec780b.png\">\r\n\r\nSigned-off-by: James-Milligan <james@omnant.co.uk>\r\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>\r\nSigned-off-by: Alex Jones <alex.jones@canonical.com>\r\nCo-authored-by: James Milligan <75740990+james-milligan@users.noreply.github.com>\r\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-06T12:45:19Z",
          "url": "https://github.com/open-feature/flagd/commit/c4e030ed651ab800aa56a555471da5913eb95259"
        },
        "date": 1663297133287,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2708,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2209809 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 8553,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "695376 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2667,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2269160 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 8294,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "722764 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2694,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2288913 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 8701,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "725793 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2594,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2325248 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 8265,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "693422 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 4021,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1514026 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 8010,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "765786 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "82278c7cf08cc6b50f49ab500caf6f9003fc0823",
          "message": "fix: upgrade package containing vulnerability (#162)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-09-16T16:56:08Z",
          "url": "https://github.com/open-feature/flagd/commit/82278c7cf08cc6b50f49ab500caf6f9003fc0823"
        },
        "date": 1663383248347,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 6289,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "906222 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2159,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2794680 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2171,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2773108 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 6262,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "921062 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2159,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2763386 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 6253,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "919413 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2147,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2802508 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 6242,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "918351 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3346,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1798567 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 6069,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "979237 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "82278c7cf08cc6b50f49ab500caf6f9003fc0823",
          "message": "fix: upgrade package containing vulnerability (#162)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-09-16T16:56:08Z",
          "url": "https://github.com/open-feature/flagd/commit/82278c7cf08cc6b50f49ab500caf6f9003fc0823"
        },
        "date": 1663469981389,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2789,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2198985 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 8992,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "586628 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2804,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2138392 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 8930,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "669986 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2810,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2136442 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 8889,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "661312 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 9031,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "647800 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2791,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2177443 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 4215,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1401628 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 8589,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "619437 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "82278c7cf08cc6b50f49ab500caf6f9003fc0823",
          "message": "fix: upgrade package containing vulnerability (#162)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-09-16T16:56:08Z",
          "url": "https://github.com/open-feature/flagd/commit/82278c7cf08cc6b50f49ab500caf6f9003fc0823"
        },
        "date": 1663556351674,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2166,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2716687 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 6462,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "924488 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2209,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2654281 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 6389,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "931954 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2175,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2768659 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 6400,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "892550 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2184,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2746844 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 6441,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "939415 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3359,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1790401 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 6109,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "935787 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "82278c7cf08cc6b50f49ab500caf6f9003fc0823",
          "message": "fix: upgrade package containing vulnerability (#162)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-09-16T16:56:08Z",
          "url": "https://github.com/open-feature/flagd/commit/82278c7cf08cc6b50f49ab500caf6f9003fc0823"
        },
        "date": 1663642594588,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2244,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2675833 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 6517,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "902511 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2245,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2681205 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 6548,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "889508 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 6473,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "888000 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2229,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2688283 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2222,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2697272 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 6473,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "900650 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3554,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1678803 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 6328,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "926689 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Michael Beemer",
            "username": "beeme1mr",
            "email": "beeme1mr@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "939b51308b89b3ed2ac057f6f7b1aac2537d56b4",
          "message": "ci: configure release please (#126)\n\n## Overview\r\n\r\nThis PR enables [Release\r\nPlease](https://github.com/googleapis/release-please) for managing\r\nautomated release PRs. It works by looking at the git history of the\r\nmain branch and automatically creating/updating a release PR.\r\nConventional commit messages are used to determine the next release\r\nversion. This repo has been configured to valid PR titles and use them\r\nas the commit message.\r\n\r\nThe flow has been tested end-to-end in a fork. You can see an example\r\nRelease Please release PR\r\n[here](https://github.com/beeme1mr/flagd/pull/7). Once approved and\r\nmerged, it triggers the release action. That can be seen\r\n[here](https://github.com/beeme1mr/flagd/actions/runs/2835311002). The\r\naction publishes the [flagD\r\ncontainer](https://github.com/beeme1mr/flagd/pkgs/container/flagd/34510492?tag=v0.1.5)\r\nand updates the [GitHub\r\nrelease](https://github.com/beeme1mr/flagd/releases/tag/v0.1.5) to\r\ninclude the binaries.\r\n\r\n## Noteworthy configurations\r\n\r\n- Release Please has been configured to **NOT** bump the major version\r\non a breaking changes until we've release a 1.0 version. The setting can\r\nbe seen\r\n[here](https://github.com/open-feature/flagd/blob/54b42caf7153133dc682c5ee91d69be5bdec218f/release-please-config.json#L7).\r\n- Goreleaser is still used but only runs after a Release Please PR has\r\nbeen merged. The setting can be seen\r\n[here](https://github.com/open-feature/flagd/blob/54b42caf7153133dc682c5ee91d69be5bdec218f/.github/workflows/release-please.yaml#L75).\r\n- Image tags had to be set manually because the action was no longer\r\ntriggered by a git tag. The setting can be seen\r\n[here](https://github.com/open-feature/flagd/blob/54b42caf7153133dc682c5ee91d69be5bdec218f/.github/workflows/release-please.yaml#L67-L69).\r\n\r\nSigned-off-by: Michael Beemer <michael.beemer@dynatrace.com>",
          "timestamp": "2022-09-20T13:01:44Z",
          "url": "https://github.com/open-feature/flagd/commit/939b51308b89b3ed2ac057f6f7b1aac2537d56b4"
        },
        "date": 1663729189448,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2449,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2434614 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 8100,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "707430 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2487,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2369137 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 7804,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "720141 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 7815,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "660956 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2484,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2350742 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 8134,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "724789 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2574,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2359147 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3894,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1554757 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 7854,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "719996 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Michael Beemer",
            "username": "beeme1mr",
            "email": "beeme1mr@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "939b51308b89b3ed2ac057f6f7b1aac2537d56b4",
          "message": "ci: configure release please (#126)\n\n## Overview\r\n\r\nThis PR enables [Release\r\nPlease](https://github.com/googleapis/release-please) for managing\r\nautomated release PRs. It works by looking at the git history of the\r\nmain branch and automatically creating/updating a release PR.\r\nConventional commit messages are used to determine the next release\r\nversion. This repo has been configured to valid PR titles and use them\r\nas the commit message.\r\n\r\nThe flow has been tested end-to-end in a fork. You can see an example\r\nRelease Please release PR\r\n[here](https://github.com/beeme1mr/flagd/pull/7). Once approved and\r\nmerged, it triggers the release action. That can be seen\r\n[here](https://github.com/beeme1mr/flagd/actions/runs/2835311002). The\r\naction publishes the [flagD\r\ncontainer](https://github.com/beeme1mr/flagd/pkgs/container/flagd/34510492?tag=v0.1.5)\r\nand updates the [GitHub\r\nrelease](https://github.com/beeme1mr/flagd/releases/tag/v0.1.5) to\r\ninclude the binaries.\r\n\r\n## Noteworthy configurations\r\n\r\n- Release Please has been configured to **NOT** bump the major version\r\non a breaking changes until we've release a 1.0 version. The setting can\r\nbe seen\r\n[here](https://github.com/open-feature/flagd/blob/54b42caf7153133dc682c5ee91d69be5bdec218f/release-please-config.json#L7).\r\n- Goreleaser is still used but only runs after a Release Please PR has\r\nbeen merged. The setting can be seen\r\n[here](https://github.com/open-feature/flagd/blob/54b42caf7153133dc682c5ee91d69be5bdec218f/.github/workflows/release-please.yaml#L75).\r\n- Image tags had to be set manually because the action was no longer\r\ntriggered by a git tag. The setting can be seen\r\n[here](https://github.com/open-feature/flagd/blob/54b42caf7153133dc682c5ee91d69be5bdec218f/.github/workflows/release-please.yaml#L67-L69).\r\n\r\nSigned-off-by: Michael Beemer <michael.beemer@dynatrace.com>",
          "timestamp": "2022-09-20T13:01:44Z",
          "url": "https://github.com/open-feature/flagd/commit/939b51308b89b3ed2ac057f6f7b1aac2537d56b4"
        },
        "date": 1663815295078,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2177,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2719318 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 6330,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "921848 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2187,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2740681 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 6346,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "911510 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2182,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2736722 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 6355,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "913147 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2169,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2767317 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 6433,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "935074 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3363,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1757758 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 6074,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "961561 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "b655e8848472679355b200f19f44b6dc55d0f0ed",
          "message": "refactor: Static reason type removal (#165)\n\nSigned-off-by: James-Milligan <james@omnant.co.uk>",
          "timestamp": "2022-09-22T13:28:53Z",
          "url": "https://github.com/open-feature/flagd/commit/b655e8848472679355b200f19f44b6dc55d0f0ed"
        },
        "date": 1663901967005,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2219,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2674090 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 6441,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "898041 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2221,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2686087 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 6478,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "899281 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2214,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2698960 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 6423,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "898406 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2216,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2707280 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 6415,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "899488 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3485,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1720825 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 6213,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "934908 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Michael Beemer",
            "username": "beeme1mr",
            "email": "beeme1mr@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "50fe46fbbf120a0657c1df35b370cdc9051d0f31",
          "message": "fix: checkout release tag before running container and binary releases (#171)\n\nSigned-off-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-23T21:10:35Z",
          "url": "https://github.com/open-feature/flagd/commit/50fe46fbbf120a0657c1df35b370cdc9051d0f31"
        },
        "date": 1663988264921,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2227,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2694912 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 6578,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "871676 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2236,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2680519 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 6564,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "886996 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2230,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2678824 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 6478,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "898248 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2223,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2693988 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 6497,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "898512 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3539,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1702285 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 6392,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "920872 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Michael Beemer",
            "username": "beeme1mr",
            "email": "beeme1mr@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "50fe46fbbf120a0657c1df35b370cdc9051d0f31",
          "message": "fix: checkout release tag before running container and binary releases (#171)\n\nSigned-off-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-23T21:10:35Z",
          "url": "https://github.com/open-feature/flagd/commit/50fe46fbbf120a0657c1df35b370cdc9051d0f31"
        },
        "date": 1664074653012,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2611,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2269186 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 7677,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "780094 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2609,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2292607 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 7558,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "763573 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2632,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2275314 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 7598,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "744343 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2603,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2279664 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 7470,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "787724 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 4000,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1494324 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 7301,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "708987 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Michael Beemer",
            "username": "beeme1mr",
            "email": "beeme1mr@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "50fe46fbbf120a0657c1df35b370cdc9051d0f31",
          "message": "fix: checkout release tag before running container and binary releases (#171)\n\nSigned-off-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-23T21:10:35Z",
          "url": "https://github.com/open-feature/flagd/commit/50fe46fbbf120a0657c1df35b370cdc9051d0f31"
        },
        "date": 1664161070597,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/happy_path",
            "value": 2210,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2714486 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveBoolean/eval_returns_error",
            "value": 6467,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "901596 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/happy_path",
            "value": 2218,
            "unit": "ns/op\t     256 B/op\t       5 allocs/op",
            "extra": "2708448 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveString/eval_returns_error",
            "value": 6493,
            "unit": "ns/op\t    1208 B/op\t      37 allocs/op",
            "extra": "901402 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/happy_path",
            "value": 2209,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2714107 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveFloat/eval_returns_error",
            "value": 6408,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "906096 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/happy_path",
            "value": 2210,
            "unit": "ns/op\t     240 B/op\t       5 allocs/op",
            "extra": "2705410 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveInt/eval_returns_error",
            "value": 6413,
            "unit": "ns/op\t    1192 B/op\t      37 allocs/op",
            "extra": "904173 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/happy_path",
            "value": 3503,
            "unit": "ns/op\t    1400 B/op\t      20 allocs/op",
            "extra": "1715283 times\n2 procs"
          },
          {
            "name": "BenchmarkGRPCService_ResolveObject/eval_returns_error",
            "value": 6221,
            "unit": "ns/op\t    1288 B/op\t      39 allocs/op",
            "extra": "941340 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "github-actions[bot]",
            "username": "github-actions[bot]",
            "email": "41898282+github-actions[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c624173bf39c53061435ff5614fc67235bdd3f38",
          "message": "chore(main): release 0.2.0 (#172)\n\n:robot: I have created a release *beep* *boop*\r\n---\r\n\r\n\r\n##\r\n[0.2.0](https://github.com/open-feature/flagd/compare/v0.1.1...v0.2.0)\r\n(2022-09-26)\r\n\r\n\r\n### ⚠ BREAKING CHANGES\r\n\r\n* Updated service to use connect (#163)\r\n\r\n### Features\r\n\r\n* Updated service to use connect\r\n([#163](https://github.com/open-feature/flagd/issues/163))\r\n([828d5c4](https://github.com/open-feature/flagd/commit/828d5c4c11157f5b7a77f5041806ba2523186764))\r\n\r\n\r\n### Bug Fixes\r\n\r\n* checkout release tag before running container and binary releases\r\n([#171](https://github.com/open-feature/flagd/issues/171))\r\n([50fe46f](https://github.com/open-feature/flagd/commit/50fe46fbbf120a0657c1df35b370cdc9051d0f31))\r\n\r\n---\r\nThis PR was generated with [Release\r\nPlease](https://github.com/googleapis/release-please). See\r\n[documentation](https://github.com/googleapis/release-please#release-please).\r\n\r\nCo-authored-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2022-09-26T13:29:13Z",
          "url": "https://github.com/open-feature/flagd/commit/c624173bf39c53061435ff5614fc67235bdd3f38"
        },
        "date": 1664247259227,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2228,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2704087 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2225,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2708293 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2196,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2716617 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2224,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2693439 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3394,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1766638 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Michael Beemer",
            "username": "beeme1mr",
            "email": "beeme1mr@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "e13caeaf02e384c3a904d0daad4b0aa951be54c3",
          "message": "ci: avoid running performance tests on forks (#173)\n\nSigned-off-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-27T18:14:26Z",
          "url": "https://github.com/open-feature/flagd/commit/e13caeaf02e384c3a904d0daad4b0aa951be54c3"
        },
        "date": 1664333782057,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2816,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2142411 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2961,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2047276 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2857,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2086869 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2839,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2108304 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4416,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1367581 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Michael Beemer",
            "username": "beeme1mr",
            "email": "beeme1mr@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "e13caeaf02e384c3a904d0daad4b0aa951be54c3",
          "message": "ci: avoid running performance tests on forks (#173)\n\nSigned-off-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-27T18:14:26Z",
          "url": "https://github.com/open-feature/flagd/commit/e13caeaf02e384c3a904d0daad4b0aa951be54c3"
        },
        "date": 1664420212035,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2310,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2603721 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2329,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2546043 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2306,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2597313 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2315,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2592352 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3615,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1666321 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Michael Beemer",
            "username": "beeme1mr",
            "email": "beeme1mr@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "e13caeaf02e384c3a904d0daad4b0aa951be54c3",
          "message": "ci: avoid running performance tests on forks (#173)\n\nSigned-off-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-27T18:14:26Z",
          "url": "https://github.com/open-feature/flagd/commit/e13caeaf02e384c3a904d0daad4b0aa951be54c3"
        },
        "date": 1664506730432,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2244,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2662776 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2228,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2677354 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2228,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2659843 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2231,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2697328 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3425,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1752847 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Michael Beemer",
            "username": "beeme1mr",
            "email": "beeme1mr@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "e13caeaf02e384c3a904d0daad4b0aa951be54c3",
          "message": "ci: avoid running performance tests on forks (#173)\n\nSigned-off-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-27T18:14:26Z",
          "url": "https://github.com/open-feature/flagd/commit/e13caeaf02e384c3a904d0daad4b0aa951be54c3"
        },
        "date": 1664593130570,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2233,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2689090 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2311,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2641231 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2234,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2682832 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2239,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2675293 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3458,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1743679 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Michael Beemer",
            "username": "beeme1mr",
            "email": "beeme1mr@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "e13caeaf02e384c3a904d0daad4b0aa951be54c3",
          "message": "ci: avoid running performance tests on forks (#173)\n\nSigned-off-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-27T18:14:26Z",
          "url": "https://github.com/open-feature/flagd/commit/e13caeaf02e384c3a904d0daad4b0aa951be54c3"
        },
        "date": 1664679546942,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2767,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2201080 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2718,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2220386 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2718,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2201997 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2673,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2186773 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4086,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1432640 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Michael Beemer",
            "username": "beeme1mr",
            "email": "beeme1mr@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "e13caeaf02e384c3a904d0daad4b0aa951be54c3",
          "message": "ci: avoid running performance tests on forks (#173)\n\nSigned-off-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-09-27T18:14:26Z",
          "url": "https://github.com/open-feature/flagd/commit/e13caeaf02e384c3a904d0daad4b0aa951be54c3"
        },
        "date": 1664765006488,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2212,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2727789 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2223,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2698213 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2220,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2708246 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2201,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2724159 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3398,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1768657 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "94d7697d08a07cede4a548ef998792d00f8954a0",
          "message": "fix: updated merge functionality (#182)\n\nSigned-off-by: James-Milligan <james@omnant.co.uk>\r\n- flag updates are handled on a per source basis, allowing deletions to\r\nbe recognised\r\n- adds the Source method to the sync provider interface, returning a\r\nstring",
          "timestamp": "2022-10-03T15:20:50Z",
          "url": "https://github.com/open-feature/flagd/commit/94d7697d08a07cede4a548ef998792d00f8954a0"
        },
        "date": 1664851641212,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2263,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2662153 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2277,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2631496 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2258,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2654494 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2264,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2658156 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3456,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1727671 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "e2fc7ee2570a119923bf95b40b2046dfa4705f20",
          "message": "feat: only fire modify event when FeatureFlagConfiguration Generation field has changed (#167)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>\r\n\r\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-04T14:44:26Z",
          "url": "https://github.com/open-feature/flagd/commit/e2fc7ee2570a119923bf95b40b2046dfa4705f20"
        },
        "date": 1664938025388,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2249,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2680092 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2244,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2678912 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2208,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2724259 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2207,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2726796 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3443,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1746374 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "e2fc7ee2570a119923bf95b40b2046dfa4705f20",
          "message": "feat: only fire modify event when FeatureFlagConfiguration Generation field has changed (#167)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>\r\n\r\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-04T14:44:26Z",
          "url": "https://github.com/open-feature/flagd/commit/e2fc7ee2570a119923bf95b40b2046dfa4705f20"
        },
        "date": 1665024356418,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2200,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2724687 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2200,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2726876 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2196,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2727628 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2226,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2735918 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3371,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1748665 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "e2fc7ee2570a119923bf95b40b2046dfa4705f20",
          "message": "feat: only fire modify event when FeatureFlagConfiguration Generation field has changed (#167)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>\r\n\r\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-04T14:44:26Z",
          "url": "https://github.com/open-feature/flagd/commit/e2fc7ee2570a119923bf95b40b2046dfa4705f20"
        },
        "date": 1665111119321,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2479,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2471450 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2452,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2452520 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2426,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2437347 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2438,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2470810 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3767,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1582912 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "e2fc7ee2570a119923bf95b40b2046dfa4705f20",
          "message": "feat: only fire modify event when FeatureFlagConfiguration Generation field has changed (#167)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>\r\n\r\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-04T14:44:26Z",
          "url": "https://github.com/open-feature/flagd/commit/e2fc7ee2570a119923bf95b40b2046dfa4705f20"
        },
        "date": 1665197061464,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2586,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2280982 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2602,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2407440 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2527,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2317875 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2493,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2361246 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3848,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1556779 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "e2fc7ee2570a119923bf95b40b2046dfa4705f20",
          "message": "feat: only fire modify event when FeatureFlagConfiguration Generation field has changed (#167)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>\r\n\r\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-04T14:44:26Z",
          "url": "https://github.com/open-feature/flagd/commit/e2fc7ee2570a119923bf95b40b2046dfa4705f20"
        },
        "date": 1665284137264,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 1747,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "3557522 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 1693,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "3575492 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 1731,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "3250612 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 1683,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "3561150 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 2619,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "2292063 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "e2fc7ee2570a119923bf95b40b2046dfa4705f20",
          "message": "feat: only fire modify event when FeatureFlagConfiguration Generation field has changed (#167)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>\r\n\r\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-04T14:44:26Z",
          "url": "https://github.com/open-feature/flagd/commit/e2fc7ee2570a119923bf95b40b2046dfa4705f20"
        },
        "date": 1665370550457,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2781,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2196042 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2780,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2172204 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2747,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2161401 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2781,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2159892 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4267,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1368502 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "3f7fcd2f57318fad4e0bf501f202af990d3c5a79",
          "message": "feat: Eventing (#187)\n\n- add eventing rpc \r\n- events are emited on successful connection to stream rpc and\r\nconfiguration changes\r\n\r\nSigned-off-by: James-Milligan <james@omnant.co.uk>\r\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-10-10T14:04:24Z",
          "url": "https://github.com/open-feature/flagd/commit/3f7fcd2f57318fad4e0bf501f202af990d3c5a79"
        },
        "date": 1665456850237,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2239,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2677734 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2261,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2661542 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2241,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2689116 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2247,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2673873 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3419,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1749813 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "6b3b6111232c132376511f4d4ea69c414f42d890",
          "message": "feat: metrics (#188)\n\nThis pull request adds prometheus metrics into the connect service using\r\na custom library ( go-http-metrics)\r\n\r\n( I decided not to go with custom gauges and counters but to use the\r\nhttp-metrics lib )\r\n\r\nAn example of the metrics we are able to leverage through\r\nhttp/https/grpc\r\n\r\n```\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"0.005\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"0.01\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"0.025\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"0.05\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"0.1\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"0.25\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"0.5\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"1\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"2.5\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"5\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"10\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"+Inf\"} 3\r\nhttp_request_duration_seconds_sum{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\"} 0.0008348739999999999\r\nhttp_request_duration_seconds_count{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\"} 3\r\n# HELP http_requests_inflight The number of inflight requests being handled at the same time.\r\n# TYPE http_requests_inflight gauge\r\nhttp_requests_inflight{handler=\"/schema.v1.Service/ResolveString\",service=\"\"} 0\r\n# HELP http_response_size_bytes The size of the HTTP responses.\r\n# TYPE http_response_size_bytes histogram\r\nhttp_response_size_bytes_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"100\"} 3\r\nhttp_response_size_bytes_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"1000\"} 3\r\nhttp_response_size_bytes_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"10000\"} 3\r\nhttp_response_size_bytes_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"100000\"} 3\r\nhttp_response_size_bytes_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"1e+06\"} 3\r\nhttp_response_size_bytes_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"1e+07\"} 3\r\nhttp_response_size_bytes_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"1e+08\"} 3\r\nhttp_response_size_bytes_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"1e+09\"} 3\r\nhttp_response_size_bytes_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"+Inf\"} 3\r\nhttp_response_size_bytes_sum{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\"} 162\r\nhttp_response_size_bytes_count{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\"} 3\r\n```\r\n\r\nThe advantage here is when this is scrapped we can perform some\r\ninteresting queries\r\n\r\n```\r\nsum(\r\n    rate(http_request_duration_seconds_count[30s])\r\n) by (handler)\r\n```\r\n\r\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>\r\nSigned-off-by: Alex Jones <alex.jones@canonical.com>\r\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-10-11T11:22:08Z",
          "url": "https://github.com/open-feature/flagd/commit/6b3b6111232c132376511f4d4ea69c414f42d890"
        },
        "date": 1665543407338,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2875,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2085986 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2845,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2088820 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2837,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2122923 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2847,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2147932 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4296,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1404974 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "6b3b6111232c132376511f4d4ea69c414f42d890",
          "message": "feat: metrics (#188)\n\nThis pull request adds prometheus metrics into the connect service using\r\na custom library ( go-http-metrics)\r\n\r\n( I decided not to go with custom gauges and counters but to use the\r\nhttp-metrics lib )\r\n\r\nAn example of the metrics we are able to leverage through\r\nhttp/https/grpc\r\n\r\n```\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"0.005\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"0.01\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"0.025\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"0.05\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"0.1\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"0.25\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"0.5\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"1\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"2.5\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"5\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"10\"} 3\r\nhttp_request_duration_seconds_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"+Inf\"} 3\r\nhttp_request_duration_seconds_sum{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\"} 0.0008348739999999999\r\nhttp_request_duration_seconds_count{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\"} 3\r\n# HELP http_requests_inflight The number of inflight requests being handled at the same time.\r\n# TYPE http_requests_inflight gauge\r\nhttp_requests_inflight{handler=\"/schema.v1.Service/ResolveString\",service=\"\"} 0\r\n# HELP http_response_size_bytes The size of the HTTP responses.\r\n# TYPE http_response_size_bytes histogram\r\nhttp_response_size_bytes_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"100\"} 3\r\nhttp_response_size_bytes_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"1000\"} 3\r\nhttp_response_size_bytes_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"10000\"} 3\r\nhttp_response_size_bytes_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"100000\"} 3\r\nhttp_response_size_bytes_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"1e+06\"} 3\r\nhttp_response_size_bytes_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"1e+07\"} 3\r\nhttp_response_size_bytes_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"1e+08\"} 3\r\nhttp_response_size_bytes_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"1e+09\"} 3\r\nhttp_response_size_bytes_bucket{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\",le=\"+Inf\"} 3\r\nhttp_response_size_bytes_sum{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\"} 162\r\nhttp_response_size_bytes_count{code=\"200\",handler=\"/schema.v1.Service/ResolveString\",method=\"POST\",service=\"\"} 3\r\n```\r\n\r\nThe advantage here is when this is scrapped we can perform some\r\ninteresting queries\r\n\r\n```\r\nsum(\r\n    rate(http_request_duration_seconds_count[30s])\r\n) by (handler)\r\n```\r\n\r\nSigned-off-by: Alex Jones <alexsimonjones@gmail.com>\r\nSigned-off-by: Alex Jones <alex.jones@canonical.com>\r\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-10-11T11:22:08Z",
          "url": "https://github.com/open-feature/flagd/commit/6b3b6111232c132376511f4d4ea69c414f42d890"
        },
        "date": 1665629764641,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2830,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2225883 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2777,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2151241 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2866,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2138178 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2787,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2116742 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4269,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1391031 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "github-actions[bot]",
            "username": "github-actions[bot]",
            "email": "41898282+github-actions[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "ccb629ebd0b4799983558a411cb7f9a4c321748f",
          "message": "chore(main): release 0.2.3 (#185)\n\n:robot: I have created a release *beep* *boop*\r\n---\r\n\r\n\r\n##\r\n[0.2.3](https://github.com/open-feature/flagd/compare/v0.2.2...v0.2.3)\r\n(2022-10-13)\r\n\r\n\r\n### Features\r\n\r\n* Eventing ([#187](https://github.com/open-feature/flagd/issues/187))\r\n([3f7fcd2](https://github.com/open-feature/flagd/commit/3f7fcd2f57318fad4e0bf501f202af990d3c5a79))\r\n* fixing informer issues\r\n([#191](https://github.com/open-feature/flagd/issues/191))\r\n([837b0c6](https://github.com/open-feature/flagd/commit/837b0c673e7e7d4799f100291ca520d22944f22a))\r\n* only fire modify event when FeatureFlagConfiguration Generation field\r\nhas changed ([#167](https://github.com/open-feature/flagd/issues/167))\r\n([e2fc7ee](https://github.com/open-feature/flagd/commit/e2fc7ee2570a119923bf95b40b2046dfa4705f20))\r\n\r\n---\r\nThis PR was generated with [Release\r\nPlease](https://github.com/googleapis/release-please). See\r\n[documentation](https://github.com/googleapis/release-please#release-please).\r\n\r\nCo-authored-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2022-10-13T11:44:19Z",
          "url": "https://github.com/open-feature/flagd/commit/ccb629ebd0b4799983558a411cb7f9a4c321748f"
        },
        "date": 1665716320253,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2362,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2556138 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2319,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2591125 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2334,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2556531 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2327,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2573079 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3555,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1686807 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "github-actions[bot]",
            "username": "github-actions[bot]",
            "email": "41898282+github-actions[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "24c8984bf76e75a6e112eeffa3809eb0e1ee9ce3",
          "message": "chore(main): release 0.2.4 (#194)\n\n:robot: I have created a release *beep* *boop*\r\n---\r\n\r\n\r\n##\r\n[0.2.4](https://github.com/open-feature/flagd/compare/v0.2.3...v0.2.4)\r\n(2022-10-14)\r\n\r\n\r\n### Bug Fixes\r\n\r\n* ApiVersion check fix\r\n([#193](https://github.com/open-feature/flagd/issues/193))\r\n([3a524d6](https://github.com/open-feature/flagd/commit/3a524d646187355bb224100f436c7b5f35abea3e))\r\n\r\n---\r\nThis PR was generated with [Release\r\nPlease](https://github.com/googleapis/release-please). See\r\n[documentation](https://github.com/googleapis/release-please#release-please).\r\n\r\nCo-authored-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2022-10-14T10:29:25Z",
          "url": "https://github.com/open-feature/flagd/commit/24c8984bf76e75a6e112eeffa3809eb0e1ee9ce3"
        },
        "date": 1665802654661,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2706,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2210160 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2712,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2223716 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2683,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2221744 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2686,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2263142 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4117,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1463990 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "github-actions[bot]",
            "username": "github-actions[bot]",
            "email": "41898282+github-actions[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "24c8984bf76e75a6e112eeffa3809eb0e1ee9ce3",
          "message": "chore(main): release 0.2.4 (#194)\n\n:robot: I have created a release *beep* *boop*\r\n---\r\n\r\n\r\n##\r\n[0.2.4](https://github.com/open-feature/flagd/compare/v0.2.3...v0.2.4)\r\n(2022-10-14)\r\n\r\n\r\n### Bug Fixes\r\n\r\n* ApiVersion check fix\r\n([#193](https://github.com/open-feature/flagd/issues/193))\r\n([3a524d6](https://github.com/open-feature/flagd/commit/3a524d646187355bb224100f436c7b5f35abea3e))\r\n\r\n---\r\nThis PR was generated with [Release\r\nPlease](https://github.com/googleapis/release-please). See\r\n[documentation](https://github.com/googleapis/release-please#release-please).\r\n\r\nCo-authored-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2022-10-14T10:29:25Z",
          "url": "https://github.com/open-feature/flagd/commit/24c8984bf76e75a6e112eeffa3809eb0e1ee9ce3"
        },
        "date": 1665889184038,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2750,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2201952 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2727,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2135546 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2708,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2236954 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2729,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2192042 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4137,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1444575 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "github-actions[bot]",
            "username": "github-actions[bot]",
            "email": "41898282+github-actions[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "24c8984bf76e75a6e112eeffa3809eb0e1ee9ce3",
          "message": "chore(main): release 0.2.4 (#194)\n\n:robot: I have created a release *beep* *boop*\r\n---\r\n\r\n\r\n##\r\n[0.2.4](https://github.com/open-feature/flagd/compare/v0.2.3...v0.2.4)\r\n(2022-10-14)\r\n\r\n\r\n### Bug Fixes\r\n\r\n* ApiVersion check fix\r\n([#193](https://github.com/open-feature/flagd/issues/193))\r\n([3a524d6](https://github.com/open-feature/flagd/commit/3a524d646187355bb224100f436c7b5f35abea3e))\r\n\r\n---\r\nThis PR was generated with [Release\r\nPlease](https://github.com/googleapis/release-please). See\r\n[documentation](https://github.com/googleapis/release-please#release-please).\r\n\r\nCo-authored-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2022-10-14T10:29:25Z",
          "url": "https://github.com/open-feature/flagd/commit/24c8984bf76e75a6e112eeffa3809eb0e1ee9ce3"
        },
        "date": 1665975557330,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2325,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2582461 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2311,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2575233 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2329,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2583627 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2321,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2580380 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3502,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1740400 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "github-actions[bot]",
            "username": "github-actions[bot]",
            "email": "41898282+github-actions[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "24c8984bf76e75a6e112eeffa3809eb0e1ee9ce3",
          "message": "chore(main): release 0.2.4 (#194)\n\n:robot: I have created a release *beep* *boop*\r\n---\r\n\r\n\r\n##\r\n[0.2.4](https://github.com/open-feature/flagd/compare/v0.2.3...v0.2.4)\r\n(2022-10-14)\r\n\r\n\r\n### Bug Fixes\r\n\r\n* ApiVersion check fix\r\n([#193](https://github.com/open-feature/flagd/issues/193))\r\n([3a524d6](https://github.com/open-feature/flagd/commit/3a524d646187355bb224100f436c7b5f35abea3e))\r\n\r\n---\r\nThis PR was generated with [Release\r\nPlease](https://github.com/googleapis/release-please). See\r\n[documentation](https://github.com/googleapis/release-please#release-please).\r\n\r\nCo-authored-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2022-10-14T10:29:25Z",
          "url": "https://github.com/open-feature/flagd/commit/24c8984bf76e75a6e112eeffa3809eb0e1ee9ce3"
        },
        "date": 1666061850891,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2336,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2572112 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2306,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2603654 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2316,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2576570 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2324,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2591156 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3525,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1690744 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "github-actions[bot]",
            "username": "github-actions[bot]",
            "email": "41898282+github-actions[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "24c8984bf76e75a6e112eeffa3809eb0e1ee9ce3",
          "message": "chore(main): release 0.2.4 (#194)\n\n:robot: I have created a release *beep* *boop*\r\n---\r\n\r\n\r\n##\r\n[0.2.4](https://github.com/open-feature/flagd/compare/v0.2.3...v0.2.4)\r\n(2022-10-14)\r\n\r\n\r\n### Bug Fixes\r\n\r\n* ApiVersion check fix\r\n([#193](https://github.com/open-feature/flagd/issues/193))\r\n([3a524d6](https://github.com/open-feature/flagd/commit/3a524d646187355bb224100f436c7b5f35abea3e))\r\n\r\n---\r\nThis PR was generated with [Release\r\nPlease](https://github.com/googleapis/release-please). See\r\n[documentation](https://github.com/googleapis/release-please#release-please).\r\n\r\nCo-authored-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2022-10-14T10:29:25Z",
          "url": "https://github.com/open-feature/flagd/commit/24c8984bf76e75a6e112eeffa3809eb0e1ee9ce3"
        },
        "date": 1666148207334,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2557,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2380012 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2450,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2504174 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2402,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2497910 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2435,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2489474 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3787,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1557800 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "github-actions[bot]",
            "username": "github-actions[bot]",
            "email": "41898282+github-actions[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "24c8984bf76e75a6e112eeffa3809eb0e1ee9ce3",
          "message": "chore(main): release 0.2.4 (#194)\n\n:robot: I have created a release *beep* *boop*\r\n---\r\n\r\n\r\n##\r\n[0.2.4](https://github.com/open-feature/flagd/compare/v0.2.3...v0.2.4)\r\n(2022-10-14)\r\n\r\n\r\n### Bug Fixes\r\n\r\n* ApiVersion check fix\r\n([#193](https://github.com/open-feature/flagd/issues/193))\r\n([3a524d6](https://github.com/open-feature/flagd/commit/3a524d646187355bb224100f436c7b5f35abea3e))\r\n\r\n---\r\nThis PR was generated with [Release\r\nPlease](https://github.com/googleapis/release-please). See\r\n[documentation](https://github.com/googleapis/release-please#release-please).\r\n\r\nCo-authored-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2022-10-14T10:29:25Z",
          "url": "https://github.com/open-feature/flagd/commit/24c8984bf76e75a6e112eeffa3809eb0e1ee9ce3"
        },
        "date": 1666234613231,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2348,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2556356 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2320,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2583052 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2335,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2572340 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2322,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2571117 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3548,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1686514 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "id": "4a9f6df4e472229ff805e9d5d3aa581c7c9c0667",
          "message": "chore: release v0.2.7\n\nRelease-As: v0.2.7",
          "timestamp": "2022-10-20T09:23:12Z",
          "url": "https://github.com/open-feature/flagd/commit/4a9f6df4e472229ff805e9d5d3aa581c7c9c0667"
        },
        "date": 1666320249796,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2344,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2557436 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2343,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2557393 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2328,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2561424 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2327,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2586153 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3562,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1693402 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1666407151378,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2859,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2111226 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2894,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2045520 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2836,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2098388 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2916,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2105436 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4443,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1344196 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1666493772013,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2300,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2636065 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2278,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2624653 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2297,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2615442 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2281,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2622702 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3387,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1769606 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1666580376191,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2313,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2592750 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2323,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2532860 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2391,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2615204 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2356,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2551267 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3537,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1697799 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1666666724067,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2970,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2051583 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2928,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2013854 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2963,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2039878 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2955,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2056857 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4506,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1352070 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1666752502316,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2823,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2115958 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2841,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2081815 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2829,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2128203 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2814,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2128275 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4419,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1354876 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1666838789479,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2348,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2553684 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2345,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2593082 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2336,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2574802 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2329,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2572238 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3553,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1687892 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1666925474501,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2333,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2519233 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2276,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2647831 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2295,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2612487 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2288,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2613732 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3417,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1753858 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1667011488263,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2310,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2580430 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2328,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2538111 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2321,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2578603 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2323,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2570316 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3536,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1696587 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1667098400193,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2334,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2574142 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2315,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2587563 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2308,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2594004 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2347,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2584590 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3531,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1701453 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1667184781989,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3025,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "1964550 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2887,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2030156 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3025,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "1989312 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3038,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "1993225 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4332,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1396528 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1667271343952,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2570,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2338122 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2474,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2145444 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2594,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2413219 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2524,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2412262 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3848,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1525194 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1667357536933,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2322,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2582937 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2293,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2574300 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2313,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2617312 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2284,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2605549 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3437,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1743292 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1667443361574,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2304,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2599149 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2292,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2624097 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2314,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2602605 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2310,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2598872 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3481,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1723921 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1667529960937,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2312,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2599756 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2299,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2572623 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2307,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2613474 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2308,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2607686 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3412,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1759922 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1667616079808,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2282,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2643159 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2238,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2693314 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2258,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2656753 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2273,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2645140 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3392,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1768978 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1667702642735,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2312,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2592016 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2287,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2590826 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2315,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2609744 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2304,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2596494 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3418,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1760085 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1667788888597,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2556,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2384269 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2703,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2252712 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2626,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2202547 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2707,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2212660 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3942,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1530717 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1667875263560,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2321,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2567919 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2322,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2567130 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2326,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2602189 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2312,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2591793 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3451,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1737279 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1667961955407,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2732,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2228594 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2672,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2210360 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2653,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2131245 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2618,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2228718 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4060,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1462030 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1668048259549,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2631,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2224941 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2561,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2286234 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2624,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2263362 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2704,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2248956 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4281,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1434334 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1668134614710,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2363,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2553540 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2336,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2558984 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2334,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2576478 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2334,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2565370 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3549,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1690351 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1668220752860,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2319,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2601674 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2303,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2623788 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2314,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2597916 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2308,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2597599 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3433,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1679638 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1668307302446,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2744,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2205056 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2733,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2160102 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2697,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2191413 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2852,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2210222 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4204,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1427990 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1668393613876,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2319,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2593294 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2306,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2598316 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2316,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2591208 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2302,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2590746 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3540,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1729800 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1668479824923,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2262,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2653238 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2245,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2686720 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2279,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2623704 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2308,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2612420 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3372,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1779034 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1668566266765,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2274,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2755909 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2167,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2719531 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2215,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2747893 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2200,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2761812 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3168,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1838214 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1668652517355,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2898,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2036726 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2809,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2055055 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2792,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2159936 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2790,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2108745 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4032,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1476003 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1668739157253,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2565,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2240587 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2571,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2330449 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2655,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2262363 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2577,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2242119 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4033,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1469492 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1668825256091,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2282,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2626520 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2297,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2603896 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2304,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2594466 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2314,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2587917 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3406,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1757586 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1668911901760,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2359,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2540263 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2331,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2516272 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2351,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2545576 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2342,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2535920 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3607,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1663536 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1668998216627,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2348,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2575204 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2312,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2599290 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2314,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2588064 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2324,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2585982 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 3538,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1699213 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "187f0f906e1c443208acf1f39026f542cbd3da2b",
          "message": "chore: contributing section for the docs (#202)\n\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-10-21T15:15:42Z",
          "url": "https://github.com/open-feature/flagd/commit/187f0f906e1c443208acf1f39026f542cbd3da2b"
        },
        "date": 1669084625009,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2845,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2133651 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2872,
            "unit": "ns/op\t     280 B/op\t       6 allocs/op",
            "extra": "2138318 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 2775,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2143119 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2835,
            "unit": "ns/op\t     264 B/op\t       6 allocs/op",
            "extra": "2148872 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4342,
            "unit": "ns/op\t    1424 B/op\t      21 allocs/op",
            "extra": "1390470 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c61984448d5cdadec62b5cf6f7e24fc5f75a3738",
          "message": "feat: snap (#211)\n\nAdding snap configuration and badge\r\n\r\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>\r\nSigned-off-by: Alex <alexsimonjones@gmail.com>\r\nCo-authored-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-11-22T21:04:42Z",
          "url": "https://github.com/open-feature/flagd/commit/c61984448d5cdadec62b5cf6f7e24fc5f75a3738"
        },
        "date": 1669170735699,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3811,
            "unit": "ns/op\t     968 B/op\t      23 allocs/op",
            "extra": "1535074 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3874,
            "unit": "ns/op\t    1000 B/op\t      24 allocs/op",
            "extra": "1551790 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 4155,
            "unit": "ns/op\t    1040 B/op\t      24 allocs/op",
            "extra": "1460250 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3827,
            "unit": "ns/op\t     968 B/op\t      23 allocs/op",
            "extra": "1560148 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5943,
            "unit": "ns/op\t    2400 B/op\t      44 allocs/op",
            "extra": "980413 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c61984448d5cdadec62b5cf6f7e24fc5f75a3738",
          "message": "feat: snap (#211)\n\nAdding snap configuration and badge\r\n\r\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>\r\nSigned-off-by: Alex <alexsimonjones@gmail.com>\r\nCo-authored-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-11-22T21:04:42Z",
          "url": "https://github.com/open-feature/flagd/commit/c61984448d5cdadec62b5cf6f7e24fc5f75a3738"
        },
        "date": 1669257037316,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3957,
            "unit": "ns/op\t     968 B/op\t      23 allocs/op",
            "extra": "1529601 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3950,
            "unit": "ns/op\t    1000 B/op\t      24 allocs/op",
            "extra": "1528051 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 4186,
            "unit": "ns/op\t    1040 B/op\t      24 allocs/op",
            "extra": "1428127 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3890,
            "unit": "ns/op\t     968 B/op\t      23 allocs/op",
            "extra": "1538034 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 6123,
            "unit": "ns/op\t    2400 B/op\t      44 allocs/op",
            "extra": "976663 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c61984448d5cdadec62b5cf6f7e24fc5f75a3738",
          "message": "feat: snap (#211)\n\nAdding snap configuration and badge\r\n\r\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>\r\nSigned-off-by: Alex <alexsimonjones@gmail.com>\r\nCo-authored-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-11-22T21:04:42Z",
          "url": "https://github.com/open-feature/flagd/commit/c61984448d5cdadec62b5cf6f7e24fc5f75a3738"
        },
        "date": 1669343494349,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3886,
            "unit": "ns/op\t     968 B/op\t      23 allocs/op",
            "extra": "1529266 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3938,
            "unit": "ns/op\t    1000 B/op\t      24 allocs/op",
            "extra": "1520931 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 4201,
            "unit": "ns/op\t    1040 B/op\t      24 allocs/op",
            "extra": "1422496 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3883,
            "unit": "ns/op\t     968 B/op\t      23 allocs/op",
            "extra": "1541649 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 6021,
            "unit": "ns/op\t    2400 B/op\t      44 allocs/op",
            "extra": "961347 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c61984448d5cdadec62b5cf6f7e24fc5f75a3738",
          "message": "feat: snap (#211)\n\nAdding snap configuration and badge\r\n\r\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>\r\nSigned-off-by: Alex <alexsimonjones@gmail.com>\r\nCo-authored-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-11-22T21:04:42Z",
          "url": "https://github.com/open-feature/flagd/commit/c61984448d5cdadec62b5cf6f7e24fc5f75a3738"
        },
        "date": 1669429762637,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 4247,
            "unit": "ns/op\t     968 B/op\t      23 allocs/op",
            "extra": "1432477 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 4218,
            "unit": "ns/op\t    1000 B/op\t      24 allocs/op",
            "extra": "1424152 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 4504,
            "unit": "ns/op\t    1040 B/op\t      24 allocs/op",
            "extra": "1331902 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 4158,
            "unit": "ns/op\t     968 B/op\t      23 allocs/op",
            "extra": "1444886 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 6145,
            "unit": "ns/op\t    2400 B/op\t      44 allocs/op",
            "extra": "945861 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c61984448d5cdadec62b5cf6f7e24fc5f75a3738",
          "message": "feat: snap (#211)\n\nAdding snap configuration and badge\r\n\r\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>\r\nSigned-off-by: Alex <alexsimonjones@gmail.com>\r\nCo-authored-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-11-22T21:04:42Z",
          "url": "https://github.com/open-feature/flagd/commit/c61984448d5cdadec62b5cf6f7e24fc5f75a3738"
        },
        "date": 1669516343820,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 4200,
            "unit": "ns/op\t     968 B/op\t      23 allocs/op",
            "extra": "1411072 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 4245,
            "unit": "ns/op\t    1000 B/op\t      24 allocs/op",
            "extra": "1412059 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 4519,
            "unit": "ns/op\t    1040 B/op\t      24 allocs/op",
            "extra": "1323824 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 4172,
            "unit": "ns/op\t     968 B/op\t      23 allocs/op",
            "extra": "1432671 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 6179,
            "unit": "ns/op\t    2400 B/op\t      44 allocs/op",
            "extra": "939331 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c61984448d5cdadec62b5cf6f7e24fc5f75a3738",
          "message": "feat: snap (#211)\n\nAdding snap configuration and badge\r\n\r\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>\r\nSigned-off-by: Alex <alexsimonjones@gmail.com>\r\nCo-authored-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-11-22T21:04:42Z",
          "url": "https://github.com/open-feature/flagd/commit/c61984448d5cdadec62b5cf6f7e24fc5f75a3738"
        },
        "date": 1669602678008,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 5431,
            "unit": "ns/op\t     968 B/op\t      23 allocs/op",
            "extra": "1000000 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 5371,
            "unit": "ns/op\t    1000 B/op\t      24 allocs/op",
            "extra": "1000000 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 5831,
            "unit": "ns/op\t    1040 B/op\t      24 allocs/op",
            "extra": "960351 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 5259,
            "unit": "ns/op\t     968 B/op\t      23 allocs/op",
            "extra": "1000000 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 8318,
            "unit": "ns/op\t    2400 B/op\t      44 allocs/op",
            "extra": "636141 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c61984448d5cdadec62b5cf6f7e24fc5f75a3738",
          "message": "feat: snap (#211)\n\nAdding snap configuration and badge\r\n\r\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>\r\nSigned-off-by: Alex <alexsimonjones@gmail.com>\r\nCo-authored-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-11-22T21:04:42Z",
          "url": "https://github.com/open-feature/flagd/commit/c61984448d5cdadec62b5cf6f7e24fc5f75a3738"
        },
        "date": 1669689011130,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 4190,
            "unit": "ns/op\t     968 B/op\t      23 allocs/op",
            "extra": "1438467 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 4226,
            "unit": "ns/op\t    1000 B/op\t      24 allocs/op",
            "extra": "1392126 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 4520,
            "unit": "ns/op\t    1040 B/op\t      24 allocs/op",
            "extra": "1326325 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 4164,
            "unit": "ns/op\t     968 B/op\t      23 allocs/op",
            "extra": "1443776 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 6168,
            "unit": "ns/op\t    2400 B/op\t      44 allocs/op",
            "extra": "946968 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Alex Jones",
            "username": "AlexsJones",
            "email": "alexsimonjones@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c61984448d5cdadec62b5cf6f7e24fc5f75a3738",
          "message": "feat: snap (#211)\n\nAdding snap configuration and badge\r\n\r\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>\r\nSigned-off-by: Alex <alexsimonjones@gmail.com>\r\nCo-authored-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-11-22T21:04:42Z",
          "url": "https://github.com/open-feature/flagd/commit/c61984448d5cdadec62b5cf6f7e24fc5f75a3738"
        },
        "date": 1669775377328,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3980,
            "unit": "ns/op\t     968 B/op\t      23 allocs/op",
            "extra": "1521369 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3984,
            "unit": "ns/op\t    1000 B/op\t      24 allocs/op",
            "extra": "1505112 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 4230,
            "unit": "ns/op\t    1040 B/op\t      24 allocs/op",
            "extra": "1412226 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3933,
            "unit": "ns/op\t     968 B/op\t      23 allocs/op",
            "extra": "1537905 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 6112,
            "unit": "ns/op\t    2400 B/op\t      44 allocs/op",
            "extra": "957856 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "14474ccf65b9b92213e8c792e94c458022484df4",
          "message": "feat!: start command flag refactor (#222)\n\nCorresponding OFO changes\r\n[here](https://github.com/open-feature/open-feature-operator/pull/256)\r\n<!-- Please use this template for your pull request. -->\r\n<!-- Please use the sections that you need and delete other sections -->\r\n\r\n## This PR\r\n<!-- add the description of the PR here -->\r\n\r\n- refactors the start command flags to remove `--sync-provider`\r\n- multiple `--uri` flags can be passed indicating the use of different\r\nexisting `sync-provider` types, all of which will work\r\n- uses a prefix on the uri to define the `sync-provider`; `http(s)://`\r\nwill be passed to the http sync, `file://` will be passed to the file\r\npath sync and the Kubernetes sync uses the following pattern\r\n`core.openfeature.dev/{namespace}/{name}`, this will also allow for the\r\nKubernetes sync to watch multiple `FeatureFlagConfigurations` from\r\ndifferent namespaces.\r\n- adds deprecation warning when the `--sync-provider` flag is passed as\r\nan argument\r\n\r\n`./flagd start --uri file://etc/flagd/end-to-end.json --uri\r\ncore.openfeature.dev/test/end-to-end-2`\r\n\r\n### Related Issues\r\n<!-- add here the GitHub issue that this PR resolves if applicable -->\r\n\r\nhttps://github.com/open-feature/open-feature-operator/issues/251\r\n\r\n### Notes\r\n<!-- any additional notes for this PR -->\r\n\r\n### Follow-up Tasks\r\n\r\n### How to test\r\n<!-- if applicable, add testing instructions under this section -->\r\n\r\nSigned-off-by: James Milligan <james@omnant.co.uk>\r\nSigned-off-by: James Milligan <75740990+james-milligan@users.noreply.github.com>\r\nCo-authored-by: Todd Baert <toddbaert@gmail.com>",
          "timestamp": "2022-11-30T21:03:10Z",
          "url": "https://github.com/open-feature/flagd/commit/14474ccf65b9b92213e8c792e94c458022484df4"
        },
        "date": 1669862150771,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 5764,
            "unit": "ns/op\t    1032 B/op\t      24 allocs/op",
            "extra": "946450 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 5971,
            "unit": "ns/op\t    1064 B/op\t      25 allocs/op",
            "extra": "972674 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 6514,
            "unit": "ns/op\t    1104 B/op\t      25 allocs/op",
            "extra": "835002 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 5837,
            "unit": "ns/op\t    1032 B/op\t      24 allocs/op",
            "extra": "994185 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 9090,
            "unit": "ns/op\t    2336 B/op\t      44 allocs/op",
            "extra": "622465 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "14474ccf65b9b92213e8c792e94c458022484df4",
          "message": "feat!: start command flag refactor (#222)\n\nCorresponding OFO changes\r\n[here](https://github.com/open-feature/open-feature-operator/pull/256)\r\n<!-- Please use this template for your pull request. -->\r\n<!-- Please use the sections that you need and delete other sections -->\r\n\r\n## This PR\r\n<!-- add the description of the PR here -->\r\n\r\n- refactors the start command flags to remove `--sync-provider`\r\n- multiple `--uri` flags can be passed indicating the use of different\r\nexisting `sync-provider` types, all of which will work\r\n- uses a prefix on the uri to define the `sync-provider`; `http(s)://`\r\nwill be passed to the http sync, `file://` will be passed to the file\r\npath sync and the Kubernetes sync uses the following pattern\r\n`core.openfeature.dev/{namespace}/{name}`, this will also allow for the\r\nKubernetes sync to watch multiple `FeatureFlagConfigurations` from\r\ndifferent namespaces.\r\n- adds deprecation warning when the `--sync-provider` flag is passed as\r\nan argument\r\n\r\n`./flagd start --uri file://etc/flagd/end-to-end.json --uri\r\ncore.openfeature.dev/test/end-to-end-2`\r\n\r\n### Related Issues\r\n<!-- add here the GitHub issue that this PR resolves if applicable -->\r\n\r\nhttps://github.com/open-feature/open-feature-operator/issues/251\r\n\r\n### Notes\r\n<!-- any additional notes for this PR -->\r\n\r\n### Follow-up Tasks\r\n\r\n### How to test\r\n<!-- if applicable, add testing instructions under this section -->\r\n\r\nSigned-off-by: James Milligan <james@omnant.co.uk>\r\nSigned-off-by: James Milligan <75740990+james-milligan@users.noreply.github.com>\r\nCo-authored-by: Todd Baert <toddbaert@gmail.com>",
          "timestamp": "2022-11-30T21:03:10Z",
          "url": "https://github.com/open-feature/flagd/commit/14474ccf65b9b92213e8c792e94c458022484df4"
        },
        "date": 1669948040469,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 4310,
            "unit": "ns/op\t    1032 B/op\t      24 allocs/op",
            "extra": "1418601 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 4207,
            "unit": "ns/op\t    1064 B/op\t      25 allocs/op",
            "extra": "1427154 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 4557,
            "unit": "ns/op\t    1104 B/op\t      25 allocs/op",
            "extra": "1314555 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 4199,
            "unit": "ns/op\t    1032 B/op\t      24 allocs/op",
            "extra": "1429946 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 6163,
            "unit": "ns/op\t    2336 B/op\t      44 allocs/op",
            "extra": "936722 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "11954b521cc6197d0dc04b163e66e38d4c288047",
          "message": "feat: enable request logging via the --debug flag (#226)\n\n<!-- Please use this template for your pull request. -->\r\n<!-- Please use the sections that you need and delete other sections -->\r\n\r\n## This PR\r\n<!-- add the description of the PR here -->\r\n\r\n- The `--debug` flag now also controls the logging for requests, when\r\nset all logging levels are enabled for the `XXXWithID` `logger` methods\r\nand allows the setting of request fields, allowing for improved\r\nperformance when logs are not required.\r\n- `NewLogger` now takes an additional boolean argument to set the\r\ninternal `reqIDLogging` field, this field is also set to false when the\r\n`*zap.Logger` argument is nil\r\n\r\n### Related Issues\r\n<!-- add here the GitHub issue that this PR resolves if applicable -->\r\n\r\n### Notes\r\n<!-- any additional notes for this PR -->\r\n\r\n### Follow-up Tasks\r\n<!-- anything that is related to this PR but not done here should be\r\nnoted under this section -->\r\n<!-- if there is a need for a new issue, please link it here -->\r\nThis flag should be set by default in the operator\r\nhttps://github.com/open-feature/open-feature-operator/pull/260\r\n### How to test\r\n<!-- if applicable, add testing instructions under this section -->\r\n\r\nSigned-off-by: James Milligan <james@omnant.co.uk>\r\nSigned-off-by: James Milligan <75740990+james-milligan@users.noreply.github.com>\r\nCo-authored-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-12-02T09:53:25Z",
          "url": "https://github.com/open-feature/flagd/commit/11954b521cc6197d0dc04b163e66e38d4c288047"
        },
        "date": 1670034022305,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3227,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1866760 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3243,
            "unit": "ns/op\t     568 B/op\t      15 allocs/op",
            "extra": "1826382 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3607,
            "unit": "ns/op\t     608 B/op\t      15 allocs/op",
            "extra": "1662622 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3230,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1858586 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5200,
            "unit": "ns/op\t    1840 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "11954b521cc6197d0dc04b163e66e38d4c288047",
          "message": "feat: enable request logging via the --debug flag (#226)\n\n<!-- Please use this template for your pull request. -->\r\n<!-- Please use the sections that you need and delete other sections -->\r\n\r\n## This PR\r\n<!-- add the description of the PR here -->\r\n\r\n- The `--debug` flag now also controls the logging for requests, when\r\nset all logging levels are enabled for the `XXXWithID` `logger` methods\r\nand allows the setting of request fields, allowing for improved\r\nperformance when logs are not required.\r\n- `NewLogger` now takes an additional boolean argument to set the\r\ninternal `reqIDLogging` field, this field is also set to false when the\r\n`*zap.Logger` argument is nil\r\n\r\n### Related Issues\r\n<!-- add here the GitHub issue that this PR resolves if applicable -->\r\n\r\n### Notes\r\n<!-- any additional notes for this PR -->\r\n\r\n### Follow-up Tasks\r\n<!-- anything that is related to this PR but not done here should be\r\nnoted under this section -->\r\n<!-- if there is a need for a new issue, please link it here -->\r\nThis flag should be set by default in the operator\r\nhttps://github.com/open-feature/open-feature-operator/pull/260\r\n### How to test\r\n<!-- if applicable, add testing instructions under this section -->\r\n\r\nSigned-off-by: James Milligan <james@omnant.co.uk>\r\nSigned-off-by: James Milligan <75740990+james-milligan@users.noreply.github.com>\r\nCo-authored-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-12-02T09:53:25Z",
          "url": "https://github.com/open-feature/flagd/commit/11954b521cc6197d0dc04b163e66e38d4c288047"
        },
        "date": 1670120744448,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3214,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1861533 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3237,
            "unit": "ns/op\t     568 B/op\t      15 allocs/op",
            "extra": "1795831 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3621,
            "unit": "ns/op\t     608 B/op\t      15 allocs/op",
            "extra": "1657676 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3228,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1857532 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5204,
            "unit": "ns/op\t    1840 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "11954b521cc6197d0dc04b163e66e38d4c288047",
          "message": "feat: enable request logging via the --debug flag (#226)\n\n<!-- Please use this template for your pull request. -->\r\n<!-- Please use the sections that you need and delete other sections -->\r\n\r\n## This PR\r\n<!-- add the description of the PR here -->\r\n\r\n- The `--debug` flag now also controls the logging for requests, when\r\nset all logging levels are enabled for the `XXXWithID` `logger` methods\r\nand allows the setting of request fields, allowing for improved\r\nperformance when logs are not required.\r\n- `NewLogger` now takes an additional boolean argument to set the\r\ninternal `reqIDLogging` field, this field is also set to false when the\r\n`*zap.Logger` argument is nil\r\n\r\n### Related Issues\r\n<!-- add here the GitHub issue that this PR resolves if applicable -->\r\n\r\n### Notes\r\n<!-- any additional notes for this PR -->\r\n\r\n### Follow-up Tasks\r\n<!-- anything that is related to this PR but not done here should be\r\nnoted under this section -->\r\n<!-- if there is a need for a new issue, please link it here -->\r\nThis flag should be set by default in the operator\r\nhttps://github.com/open-feature/open-feature-operator/pull/260\r\n### How to test\r\n<!-- if applicable, add testing instructions under this section -->\r\n\r\nSigned-off-by: James Milligan <james@omnant.co.uk>\r\nSigned-off-by: James Milligan <75740990+james-milligan@users.noreply.github.com>\r\nCo-authored-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2022-12-02T09:53:25Z",
          "url": "https://github.com/open-feature/flagd/commit/11954b521cc6197d0dc04b163e66e38d4c288047"
        },
        "date": 1670207088864,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3009,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "2011780 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3041,
            "unit": "ns/op\t     568 B/op\t      15 allocs/op",
            "extra": "1977130 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3327,
            "unit": "ns/op\t     608 B/op\t      15 allocs/op",
            "extra": "1812774 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3021,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1981984 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4798,
            "unit": "ns/op\t    1840 B/op\t      34 allocs/op",
            "extra": "1248081 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "github-actions[bot]",
            "username": "github-actions[bot]",
            "email": "41898282+github-actions[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "df7c6ee15ab30325b893cdce95e0735d486ebf74",
          "message": "chore(main): release 0.2.7 (#215)\n\n:robot: I have created a release *beep* *boop*\r\n---\r\n\r\n\r\n##\r\n[0.2.7](https://github.com/open-feature/flagd/compare/v0.2.5...v0.2.7)\r\n(2022-12-02)\r\n\r\n\r\n### ⚠ BREAKING CHANGES\r\n\r\n* start command flag refactor\r\n([#222](https://github.com/open-feature/flagd/issues/222))\r\n\r\n### Features\r\n\r\n* enable request logging via the --debug flag\r\n([#226](https://github.com/open-feature/flagd/issues/226))\r\n([11954b5](https://github.com/open-feature/flagd/commit/11954b521cc6197d0dc04b163e66e38d4c288047))\r\n* Resurrected the STATIC flag reason. Documented the caching strategy.\r\n([#224](https://github.com/open-feature/flagd/issues/224))\r\n([5830592](https://github.com/open-feature/flagd/commit/5830592053c55dc9e55c16854e40c3fc8345d6d1))\r\n* snap ([#211](https://github.com/open-feature/flagd/issues/211))\r\n([c619844](https://github.com/open-feature/flagd/commit/c61984448d5cdadec62b5cf6f7e24fc5f75a3738))\r\n* start command flag refactor\r\n([#222](https://github.com/open-feature/flagd/issues/222))\r\n([14474cc](https://github.com/open-feature/flagd/commit/14474ccf65b9b92213e8c792e94c458022484df4))\r\n\r\n\r\n### Miscellaneous Chores\r\n\r\n* release v0.2.6\r\n([93cfb78](https://github.com/open-feature/flagd/commit/93cfb78d024b436fa7fb17fd41f74d1508bf8b64))\r\n* release v0.2.7\r\n([4a9f6df](https://github.com/open-feature/flagd/commit/4a9f6df4e472229ff805e9d5d3aa581c7c9c0667))\r\n\r\n---\r\nThis PR was generated with [Release\r\nPlease](https://github.com/googleapis/release-please). See\r\n[documentation](https://github.com/googleapis/release-please#release-please).\r\n\r\nCo-authored-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2022-12-05T12:47:44Z",
          "url": "https://github.com/open-feature/flagd/commit/df7c6ee15ab30325b893cdce95e0735d486ebf74"
        },
        "date": 1670293510003,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3057,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1994467 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3040,
            "unit": "ns/op\t     568 B/op\t      15 allocs/op",
            "extra": "1967996 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3353,
            "unit": "ns/op\t     608 B/op\t      15 allocs/op",
            "extra": "1785908 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3069,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1962642 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4924,
            "unit": "ns/op\t    1840 B/op\t      34 allocs/op",
            "extra": "1235658 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "5bbef9ea4b1960686e58298c2c2e192ca99f072f",
          "message": "fix: changed eventing configuration mutex to rwmutex and added missing lock (#220)\n\nFixes https://github.com/open-feature/flagd/issues/219",
          "timestamp": "2022-12-06T15:12:20Z",
          "url": "https://github.com/open-feature/flagd/commit/5bbef9ea4b1960686e58298c2c2e192ca99f072f"
        },
        "date": 1670380187182,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3649,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1608765 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3663,
            "unit": "ns/op\t     568 B/op\t      15 allocs/op",
            "extra": "1643257 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 4038,
            "unit": "ns/op\t     608 B/op\t      15 allocs/op",
            "extra": "1484007 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3649,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1634730 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 6036,
            "unit": "ns/op\t    1840 B/op\t      34 allocs/op",
            "extra": "951963 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "5bbef9ea4b1960686e58298c2c2e192ca99f072f",
          "message": "fix: changed eventing configuration mutex to rwmutex and added missing lock (#220)\n\nFixes https://github.com/open-feature/flagd/issues/219",
          "timestamp": "2022-12-06T15:12:20Z",
          "url": "https://github.com/open-feature/flagd/commit/5bbef9ea4b1960686e58298c2c2e192ca99f072f"
        },
        "date": 1670466396811,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3053,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1964342 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3085,
            "unit": "ns/op\t     568 B/op\t      15 allocs/op",
            "extra": "1903726 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3363,
            "unit": "ns/op\t     608 B/op\t      15 allocs/op",
            "extra": "1782285 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3070,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1942149 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4897,
            "unit": "ns/op\t    1840 B/op\t      34 allocs/op",
            "extra": "1227698 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "5bbef9ea4b1960686e58298c2c2e192ca99f072f",
          "message": "fix: changed eventing configuration mutex to rwmutex and added missing lock (#220)\n\nFixes https://github.com/open-feature/flagd/issues/219",
          "timestamp": "2022-12-06T15:12:20Z",
          "url": "https://github.com/open-feature/flagd/commit/5bbef9ea4b1960686e58298c2c2e192ca99f072f"
        },
        "date": 1670552789270,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3113,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1959466 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3066,
            "unit": "ns/op\t     568 B/op\t      15 allocs/op",
            "extra": "1961010 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3417,
            "unit": "ns/op\t     608 B/op\t      15 allocs/op",
            "extra": "1783951 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3087,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1881186 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5044,
            "unit": "ns/op\t    1840 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "d6147490fd939f227025b2875b2a1888b4e4a3fe",
          "message": "docs: Restructure flagd docs (#231)\n\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-12-09T18:52:30Z",
          "url": "https://github.com/open-feature/flagd/commit/d6147490fd939f227025b2875b2a1888b4e4a3fe"
        },
        "date": 1670638988949,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2992,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1970922 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3025,
            "unit": "ns/op\t     568 B/op\t      15 allocs/op",
            "extra": "1986721 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3291,
            "unit": "ns/op\t     608 B/op\t      15 allocs/op",
            "extra": "1814101 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3000,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "2003631 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4815,
            "unit": "ns/op\t    1840 B/op\t      34 allocs/op",
            "extra": "1248463 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "d6147490fd939f227025b2875b2a1888b4e4a3fe",
          "message": "docs: Restructure flagd docs (#231)\n\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-12-09T18:52:30Z",
          "url": "https://github.com/open-feature/flagd/commit/d6147490fd939f227025b2875b2a1888b4e4a3fe"
        },
        "date": 1670725763906,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3071,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1981899 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3071,
            "unit": "ns/op\t     568 B/op\t      15 allocs/op",
            "extra": "1950528 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3333,
            "unit": "ns/op\t     608 B/op\t      15 allocs/op",
            "extra": "1790776 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3000,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1989630 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4831,
            "unit": "ns/op\t    1840 B/op\t      34 allocs/op",
            "extra": "1244931 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "d6147490fd939f227025b2875b2a1888b4e4a3fe",
          "message": "docs: Restructure flagd docs (#231)\n\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-12-09T18:52:30Z",
          "url": "https://github.com/open-feature/flagd/commit/d6147490fd939f227025b2875b2a1888b4e4a3fe"
        },
        "date": 1670812031095,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3593,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1633480 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3687,
            "unit": "ns/op\t     568 B/op\t      15 allocs/op",
            "extra": "1639299 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3998,
            "unit": "ns/op\t     608 B/op\t      15 allocs/op",
            "extra": "1509178 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3658,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1658766 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5837,
            "unit": "ns/op\t    1840 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Michael Beemer",
            "username": "beeme1mr",
            "email": "beeme1mr@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "c2a9ac9c72e775b0da2e055c850d5288b0d2708c",
          "message": "chore: update architecture doc image path\n\nSigned-off-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2022-12-12T19:00:19Z",
          "url": "https://github.com/open-feature/flagd/commit/c2a9ac9c72e775b0da2e055c850d5288b0d2708c"
        },
        "date": 1670898578923,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3238,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1794806 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3321,
            "unit": "ns/op\t     568 B/op\t      15 allocs/op",
            "extra": "1778922 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3601,
            "unit": "ns/op\t     608 B/op\t      15 allocs/op",
            "extra": "1632300 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3257,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1825480 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5408,
            "unit": "ns/op\t    1840 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Kavindu Dodanduwa",
            "username": "Kavindu-Dodan",
            "email": "Kavindu-Dodan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "26de63e3ced390f28674c624f3ee5438a4641843",
          "message": "test: load testing tooling and initial observations (#225)\n\nSigned-off-by: Kavindu Dodanduwa <KavinduDodanduwa@gmail.com>",
          "timestamp": "2022-12-13T12:41:21Z",
          "url": "https://github.com/open-feature/flagd/commit/26de63e3ced390f28674c624f3ee5438a4641843"
        },
        "date": 1670985062745,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11926,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "499651 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1177,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5114890 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1326,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4511864 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1332,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4515289 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11864,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "501033 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1172,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5108000 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1326,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4525882 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1318,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4552922 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11926,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "495187 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1150,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5209148 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1329,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4523948 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1329,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4513557 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11032,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "538988 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1156,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5204713 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1326,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4529054 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1320,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4542026 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1162,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5168990 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1329,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4528800 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1327,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4517617 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3230,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1863574 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3243,
            "unit": "ns/op\t     568 B/op\t      15 allocs/op",
            "extra": "1848230 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3623,
            "unit": "ns/op\t     608 B/op\t      15 allocs/op",
            "extra": "1657288 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3233,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1858371 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5182,
            "unit": "ns/op\t    1840 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Kavindu Dodanduwa",
            "username": "Kavindu-Dodan",
            "email": "Kavindu-Dodan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "26de63e3ced390f28674c624f3ee5438a4641843",
          "message": "test: load testing tooling and initial observations (#225)\n\nSigned-off-by: Kavindu Dodanduwa <KavinduDodanduwa@gmail.com>",
          "timestamp": "2022-12-13T12:41:21Z",
          "url": "https://github.com/open-feature/flagd/commit/26de63e3ced390f28674c624f3ee5438a4641843"
        },
        "date": 1671071438816,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11526,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "514468 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1140,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5275963 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1292,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4652832 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1287,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4685728 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11460,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "512088 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1146,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5239179 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1300,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4634065 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1297,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4652624 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11553,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "483409 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1138,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5285320 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1297,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4606606 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1300,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4625827 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 10507,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "504802 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1135,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5264130 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1286,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4655295 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1280,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4680595 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1164,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5142340 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1305,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4602561 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1302,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4631728 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3037,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1999206 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3019,
            "unit": "ns/op\t     568 B/op\t      15 allocs/op",
            "extra": "2006098 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3279,
            "unit": "ns/op\t     608 B/op\t      15 allocs/op",
            "extra": "1831424 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2971,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "2022406 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4773,
            "unit": "ns/op\t    1840 B/op\t      34 allocs/op",
            "extra": "1258005 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Kavindu Dodanduwa",
            "username": "Kavindu-Dodan",
            "email": "Kavindu-Dodan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "26de63e3ced390f28674c624f3ee5438a4641843",
          "message": "test: load testing tooling and initial observations (#225)\n\nSigned-off-by: Kavindu Dodanduwa <KavinduDodanduwa@gmail.com>",
          "timestamp": "2022-12-13T12:41:21Z",
          "url": "https://github.com/open-feature/flagd/commit/26de63e3ced390f28674c624f3ee5438a4641843"
        },
        "date": 1671157263071,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11491,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "515749 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1126,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5318418 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1288,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4653206 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1290,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4669402 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11469,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "513620 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1156,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5198358 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1308,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4587606 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1301,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4595164 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11579,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "496748 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1146,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5241373 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1300,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4610544 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1292,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4667184 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 10535,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "557517 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1134,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5294788 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1293,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4652578 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1297,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4629890 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1153,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5181927 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1301,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4605015 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1296,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4631793 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3025,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1998313 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3074,
            "unit": "ns/op\t     568 B/op\t      15 allocs/op",
            "extra": "1970193 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3340,
            "unit": "ns/op\t     608 B/op\t      15 allocs/op",
            "extra": "1791226 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3044,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1965829 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4853,
            "unit": "ns/op\t    1840 B/op\t      34 allocs/op",
            "extra": "1249836 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Kavindu Dodanduwa",
            "username": "Kavindu-Dodan",
            "email": "Kavindu-Dodan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "26de63e3ced390f28674c624f3ee5438a4641843",
          "message": "test: load testing tooling and initial observations (#225)\n\nSigned-off-by: Kavindu Dodanduwa <KavinduDodanduwa@gmail.com>",
          "timestamp": "2022-12-13T12:41:21Z",
          "url": "https://github.com/open-feature/flagd/commit/26de63e3ced390f28674c624f3ee5438a4641843"
        },
        "date": 1671243628181,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 12027,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "494746 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1160,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5173383 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1332,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4505845 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1324,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4529554 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11914,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "497766 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1172,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5091132 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1334,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4494306 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1324,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4539198 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11934,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "491334 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1166,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5132869 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1335,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4503201 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1331,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4503912 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11045,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "534615 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1164,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5160228 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1327,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4500448 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1334,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4482832 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1170,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5119078 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1336,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4500420 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1329,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4519930 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3205,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1871179 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3250,
            "unit": "ns/op\t     568 B/op\t      15 allocs/op",
            "extra": "1812246 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3588,
            "unit": "ns/op\t     608 B/op\t      15 allocs/op",
            "extra": "1669188 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3231,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1857320 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5152,
            "unit": "ns/op\t    1840 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Kavindu Dodanduwa",
            "username": "Kavindu-Dodan",
            "email": "Kavindu-Dodan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "26de63e3ced390f28674c624f3ee5438a4641843",
          "message": "test: load testing tooling and initial observations (#225)\n\nSigned-off-by: Kavindu Dodanduwa <KavinduDodanduwa@gmail.com>",
          "timestamp": "2022-12-13T12:41:21Z",
          "url": "https://github.com/open-feature/flagd/commit/26de63e3ced390f28674c624f3ee5438a4641843"
        },
        "date": 1671330141754,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 12037,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "487555 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1148,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5228048 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1297,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4622272 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1304,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4578894 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11799,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "485964 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1157,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5204326 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1310,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4530945 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1302,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4621058 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11943,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "502749 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1146,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5231128 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1311,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4524847 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1300,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4629206 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 10804,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "544563 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1151,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5211726 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1293,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4625331 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1297,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4614616 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1147,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5233440 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1307,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4620168 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1299,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4614054 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2973,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1954773 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3044,
            "unit": "ns/op\t     568 B/op\t      15 allocs/op",
            "extra": "1965848 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3312,
            "unit": "ns/op\t     608 B/op\t      15 allocs/op",
            "extra": "1819797 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3021,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1998132 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4884,
            "unit": "ns/op\t    1840 B/op\t      34 allocs/op",
            "extra": "1237700 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Kavindu Dodanduwa",
            "username": "Kavindu-Dodan",
            "email": "Kavindu-Dodan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "26de63e3ced390f28674c624f3ee5438a4641843",
          "message": "test: load testing tooling and initial observations (#225)\n\nSigned-off-by: Kavindu Dodanduwa <KavinduDodanduwa@gmail.com>",
          "timestamp": "2022-12-13T12:41:21Z",
          "url": "https://github.com/open-feature/flagd/commit/26de63e3ced390f28674c624f3ee5438a4641843"
        },
        "date": 1671416448336,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11605,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "502291 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1140,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5232338 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1299,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4591106 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1280,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4667752 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11529,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "510278 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1156,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5178254 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1301,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4572186 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1302,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4638414 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11669,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "500218 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1151,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5214175 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1310,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4579509 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1300,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4619886 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 10677,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "548840 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1151,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5199982 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1291,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4639003 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1286,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4676187 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1143,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5248390 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1312,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4566518 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1292,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4653334 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2943,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "2029692 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3021,
            "unit": "ns/op\t     568 B/op\t      15 allocs/op",
            "extra": "1961864 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3281,
            "unit": "ns/op\t     608 B/op\t      15 allocs/op",
            "extra": "1831873 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2987,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "2010782 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4780,
            "unit": "ns/op\t    1840 B/op\t      34 allocs/op",
            "extra": "1260580 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Kavindu Dodanduwa",
            "username": "Kavindu-Dodan",
            "email": "Kavindu-Dodan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "26de63e3ced390f28674c624f3ee5438a4641843",
          "message": "test: load testing tooling and initial observations (#225)\n\nSigned-off-by: Kavindu Dodanduwa <KavinduDodanduwa@gmail.com>",
          "timestamp": "2022-12-13T12:41:21Z",
          "url": "https://github.com/open-feature/flagd/commit/26de63e3ced390f28674c624f3ee5438a4641843"
        },
        "date": 1671503086528,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 12026,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "490689 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1160,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5154447 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1330,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4502108 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1332,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4480510 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11952,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "496939 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1168,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5128686 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1341,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4506082 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1325,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4529983 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 12086,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "492804 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1156,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5197941 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1335,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4499227 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1323,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4544894 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11109,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "532177 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1154,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5190820 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1330,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4512373 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1324,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4529223 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1165,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5145458 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1330,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4501303 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1330,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4508499 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3199,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1873321 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3269,
            "unit": "ns/op\t     568 B/op\t      15 allocs/op",
            "extra": "1799700 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3596,
            "unit": "ns/op\t     608 B/op\t      15 allocs/op",
            "extra": "1668139 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3257,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1845160 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5150,
            "unit": "ns/op\t    1840 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Kavindu Dodanduwa",
            "username": "Kavindu-Dodan",
            "email": "Kavindu-Dodan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "26de63e3ced390f28674c624f3ee5438a4641843",
          "message": "test: load testing tooling and initial observations (#225)\n\nSigned-off-by: Kavindu Dodanduwa <KavinduDodanduwa@gmail.com>",
          "timestamp": "2022-12-13T12:41:21Z",
          "url": "https://github.com/open-feature/flagd/commit/26de63e3ced390f28674c624f3ee5438a4641843"
        },
        "date": 1671589326384,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11603,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "494148 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1142,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5271643 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1305,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4573518 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1297,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4644025 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11680,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "503774 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1143,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5209842 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1310,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4571673 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1291,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4681708 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11810,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "494761 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1157,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5157663 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1315,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4548306 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1300,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4608912 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 10774,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "538236 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1139,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5275779 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1303,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4560612 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1295,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4634107 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1147,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5209444 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1312,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4625636 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1313,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4564764 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3045,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1988316 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3059,
            "unit": "ns/op\t     568 B/op\t      15 allocs/op",
            "extra": "1908654 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3379,
            "unit": "ns/op\t     608 B/op\t      15 allocs/op",
            "extra": "1776768 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3063,
            "unit": "ns/op\t     536 B/op\t      14 allocs/op",
            "extra": "1962178 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4818,
            "unit": "ns/op\t    1840 B/op\t      34 allocs/op",
            "extra": "1247695 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "af6f133688b506f52ac6180890e3d785238d7b40",
          "message": "chore: buf upgrade to v1 (#229)",
          "timestamp": "2022-12-21T12:28:07Z",
          "url": "https://github.com/open-feature/flagd/commit/af6f133688b506f52ac6180890e3d785238d7b40"
        },
        "date": 1671675813453,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 12004,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "493626 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1150,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5178752 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1330,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4511475 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1338,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4475786 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11899,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "499728 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1166,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5151334 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1327,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4526750 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1319,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4546278 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11980,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "487923 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1156,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5192214 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1325,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4509399 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1322,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4530273 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11077,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "532538 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1149,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5202933 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1330,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4488790 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1327,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4532295 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1167,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5133537 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1334,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4505160 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1346,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4496468 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3091,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1942611 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3113,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1876411 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3514,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1706128 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3178,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1891210 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5201,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "af6f133688b506f52ac6180890e3d785238d7b40",
          "message": "chore: buf upgrade to v1 (#229)",
          "timestamp": "2022-12-21T12:28:07Z",
          "url": "https://github.com/open-feature/flagd/commit/af6f133688b506f52ac6180890e3d785238d7b40"
        },
        "date": 1671762232390,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 12002,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "495288 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1153,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5197734 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1330,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4519264 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1332,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4501791 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11974,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "496774 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1158,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5164203 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1331,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4507378 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1337,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4488615 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11966,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "493682 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1162,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5167278 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1333,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4498932 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1325,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4517790 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11036,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "528625 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1166,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5170006 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1343,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4491757 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1323,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4541133 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1172,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5125929 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1336,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4486621 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1337,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4485470 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3147,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1937007 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3115,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1922846 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3546,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1696627 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3198,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1877761 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5134,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "af6f133688b506f52ac6180890e3d785238d7b40",
          "message": "chore: buf upgrade to v1 (#229)",
          "timestamp": "2022-12-21T12:28:07Z",
          "url": "https://github.com/open-feature/flagd/commit/af6f133688b506f52ac6180890e3d785238d7b40"
        },
        "date": 1671848357257,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11974,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "494636 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1156,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5154672 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1327,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4535712 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1323,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4523901 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11889,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "496881 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1165,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5145505 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1332,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4506535 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1327,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4518775 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11915,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "495836 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1175,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5106895 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1325,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4530769 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1322,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4541380 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11013,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "538532 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1167,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5156683 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1328,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4492575 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1321,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4530560 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1175,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5089503 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1324,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4531586 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1330,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4507239 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3087,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1942317 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3117,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1892346 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3550,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1708262 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3167,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1896572 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5095,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "af6f133688b506f52ac6180890e3d785238d7b40",
          "message": "chore: buf upgrade to v1 (#229)",
          "timestamp": "2022-12-21T12:28:07Z",
          "url": "https://github.com/open-feature/flagd/commit/af6f133688b506f52ac6180890e3d785238d7b40"
        },
        "date": 1671935126556,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11651,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "503331 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1154,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5173387 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1291,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4651928 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1293,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4644000 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11674,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "512361 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1147,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5237396 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1303,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4587896 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1305,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4609662 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11785,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "483108 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1147,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5229189 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1309,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4559666 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1307,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4589311 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 10843,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "544042 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1141,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5255253 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1296,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4578450 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1287,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4655544 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1152,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5221947 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1302,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4599343 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1302,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4607570 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2995,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2047993 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2942,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "2038326 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3296,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1824802 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2947,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2036018 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4835,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1235430 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "af6f133688b506f52ac6180890e3d785238d7b40",
          "message": "chore: buf upgrade to v1 (#229)",
          "timestamp": "2022-12-21T12:28:07Z",
          "url": "https://github.com/open-feature/flagd/commit/af6f133688b506f52ac6180890e3d785238d7b40"
        },
        "date": 1672021458002,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 12617,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "466508 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1321,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "4684434 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1514,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "3893928 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1509,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4022815 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 12543,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "468039 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1318,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4590384 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1536,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3867555 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1522,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3965253 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 12847,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "454314 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1318,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4735184 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1507,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4045170 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1485,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4030546 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11422,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "515880 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1308,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4519132 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1521,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4068949 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1505,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4044452 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1332,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4626428 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1513,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3949410 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1527,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3851272 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3271,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1852143 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3306,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1824690 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3712,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1634798 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3253,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1787548 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5290,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "af6f133688b506f52ac6180890e3d785238d7b40",
          "message": "chore: buf upgrade to v1 (#229)",
          "timestamp": "2022-12-21T12:28:07Z",
          "url": "https://github.com/open-feature/flagd/commit/af6f133688b506f52ac6180890e3d785238d7b40"
        },
        "date": 1672107716208,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 12013,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "495502 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1146,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5232440 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1338,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4520864 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1325,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4536326 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11881,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "496904 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1167,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5149404 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1329,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4494177 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1324,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4519440 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11934,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "489433 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1156,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5184163 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1332,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4518776 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1331,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4508426 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 10991,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "535256 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1163,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4946901 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1346,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4458434 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1334,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4487980 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1178,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5087920 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1322,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4531498 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1328,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4519017 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3217,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1876411 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3223,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1864060 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3588,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1671175 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3218,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1862498 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5171,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "af6f133688b506f52ac6180890e3d785238d7b40",
          "message": "chore: buf upgrade to v1 (#229)",
          "timestamp": "2022-12-21T12:28:07Z",
          "url": "https://github.com/open-feature/flagd/commit/af6f133688b506f52ac6180890e3d785238d7b40"
        },
        "date": 1672194188486,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11928,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "502896 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1166,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5132767 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1304,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4580005 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1327,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4513054 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11796,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "499690 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1166,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5151747 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1321,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4561706 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1314,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4552707 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11906,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "483702 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1157,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5188144 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1316,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4545378 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1313,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4579754 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 10804,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "537568 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1138,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5268705 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1303,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4611028 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1310,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4577232 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1148,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5225082 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1321,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4537488 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1309,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4573189 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2999,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2019778 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2964,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "2018853 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3289,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1829347 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2964,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2019570 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4814,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1250820 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "af6f133688b506f52ac6180890e3d785238d7b40",
          "message": "chore: buf upgrade to v1 (#229)",
          "timestamp": "2022-12-21T12:28:07Z",
          "url": "https://github.com/open-feature/flagd/commit/af6f133688b506f52ac6180890e3d785238d7b40"
        },
        "date": 1672280656458,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11435,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "515905 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1144,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5247346 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1289,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4656565 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1286,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4646534 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11404,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "516034 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1143,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5251322 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1298,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4611886 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1285,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4664071 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11503,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "509376 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1146,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5224222 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1298,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4639868 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1287,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4651533 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 10468,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "560301 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1136,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5296261 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1282,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4668560 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1300,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4643194 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1149,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5219229 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1297,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4657088 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1298,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4616745 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2906,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2052199 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2906,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "2047692 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3252,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1848519 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2887,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2080102 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4701,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1272616 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "af6f133688b506f52ac6180890e3d785238d7b40",
          "message": "chore: buf upgrade to v1 (#229)",
          "timestamp": "2022-12-21T12:28:07Z",
          "url": "https://github.com/open-feature/flagd/commit/af6f133688b506f52ac6180890e3d785238d7b40"
        },
        "date": 1672367055953,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11436,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "511352 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1148,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5231941 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1297,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4633207 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1287,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4667044 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11390,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "511612 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1132,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5281130 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1296,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4638354 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1304,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4589689 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11547,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "507290 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1157,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5182262 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1300,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4611942 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1298,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4623021 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 10558,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "551641 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1146,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5236275 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1298,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4608100 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1290,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4683076 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1152,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5211375 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1312,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4579095 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1300,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4621810 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2944,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2053569 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2953,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "2024472 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3279,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1842531 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2979,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2026136 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4822,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1241034 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "af6f133688b506f52ac6180890e3d785238d7b40",
          "message": "chore: buf upgrade to v1 (#229)",
          "timestamp": "2022-12-21T12:28:07Z",
          "url": "https://github.com/open-feature/flagd/commit/af6f133688b506f52ac6180890e3d785238d7b40"
        },
        "date": 1672453308600,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11785,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "498207 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1131,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5270086 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1301,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4622067 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1294,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4650414 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11820,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "497365 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1149,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5229813 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1304,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4605288 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1308,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4582128 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 12042,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "492429 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1153,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5180296 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1315,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4564201 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1307,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4566009 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 10846,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "533778 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1138,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5269959 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1293,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4646384 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1305,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4606665 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1148,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5168349 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1317,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4557056 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1305,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4597528 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2982,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1974187 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2985,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "2007426 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3317,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1807143 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2984,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2005984 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4863,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1230948 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "af6f133688b506f52ac6180890e3d785238d7b40",
          "message": "chore: buf upgrade to v1 (#229)",
          "timestamp": "2022-12-21T12:28:07Z",
          "url": "https://github.com/open-feature/flagd/commit/af6f133688b506f52ac6180890e3d785238d7b40"
        },
        "date": 1672540410989,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 13843,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "420286 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1379,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "4303918 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1597,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "3781406 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1589,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "3791494 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 13714,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "432884 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1411,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4233690 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1621,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3682958 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1594,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3790848 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 13880,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "435474 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1398,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4256528 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1595,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3724684 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1609,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3732720 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 12409,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "483028 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1397,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4292563 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1576,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "3766852 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1584,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3807171 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1407,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4311208 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1576,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3647019 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1607,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3793987 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3544,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1727227 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3505,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1702291 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3941,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1530592 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3539,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1706348 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5682,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "af6f133688b506f52ac6180890e3d785238d7b40",
          "message": "chore: buf upgrade to v1 (#229)",
          "timestamp": "2022-12-21T12:28:07Z",
          "url": "https://github.com/open-feature/flagd/commit/af6f133688b506f52ac6180890e3d785238d7b40"
        },
        "date": 1672626248567,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 12241,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "485869 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1165,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5156004 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1341,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4460412 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1355,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4482104 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 12224,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "487074 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1157,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5209478 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1348,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4429387 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1342,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4449604 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 12146,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "484458 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1180,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5092905 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1341,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4406749 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1366,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4445188 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11427,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "520491 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1173,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5148086 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1345,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4471227 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1335,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4492972 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1186,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5066642 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1342,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4452354 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1343,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4452252 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3169,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1920728 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3196,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1903603 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3565,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1678269 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3233,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1854910 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5241,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "James Milligan",
            "username": "james-milligan",
            "email": "75740990+james-milligan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "af6f133688b506f52ac6180890e3d785238d7b40",
          "message": "chore: buf upgrade to v1 (#229)",
          "timestamp": "2022-12-21T12:28:07Z",
          "url": "https://github.com/open-feature/flagd/commit/af6f133688b506f52ac6180890e3d785238d7b40"
        },
        "date": 1672712586965,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11502,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "511604 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1129,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5273205 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1286,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4685746 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1289,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4676149 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11485,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "511120 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1137,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5292342 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1293,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4629452 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1298,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4620045 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11587,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "506679 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1139,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5267438 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1299,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4626099 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1286,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4680390 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 10540,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "548224 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1139,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5266021 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1281,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4673725 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1286,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4681402 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1136,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5283666 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1300,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4622655 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1280,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4669759 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2919,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2057738 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2918,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "2006353 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3214,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1858738 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2926,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2049685 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4732,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1266672 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Todd Baert",
            "username": "toddbaert",
            "email": "toddbaert@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "3ecfc654f92cff7ce6ba7f802d5196202f88c653",
          "message": "chore: update CODEOWNERS (#228)\n\nCo-authored-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2023-01-03T20:29:38Z",
          "url": "https://github.com/open-feature/flagd/commit/3ecfc654f92cff7ce6ba7f802d5196202f88c653"
        },
        "date": 1672799185823,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 12065,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "490000 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1146,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5231775 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1322,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4495669 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1334,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4494880 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 12153,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "492648 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1150,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5201827 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1329,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4521482 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1334,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4457710 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11929,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "495256 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1155,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5200920 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1325,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4524048 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1325,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4519353 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11031,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "534572 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1147,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5228133 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1327,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4521879 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1312,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4567821 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1170,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5108085 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1324,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4527206 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1324,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4526275 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3077,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1896390 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3102,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1933491 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3514,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1708502 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3163,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1899526 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5089,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Suraj Banakar(बानकर) | スラジ",
            "username": "vadasambar",
            "email": "34534103+vadasambar@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "2dbace5b6bb8e187a7d44a3d3ec14190c63b3ae0",
          "message": "feat: support yaml evaluator (#206)\n\nSigned-off-by: Suraj Banakar <surajrbanakar@gmail.com>",
          "timestamp": "2023-01-04T13:26:16Z",
          "url": "https://github.com/open-feature/flagd/commit/2dbace5b6bb8e187a7d44a3d3ec14190c63b3ae0"
        },
        "date": 1672885608788,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11951,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "496908 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1159,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5162269 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1330,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4509055 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1324,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4531190 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11907,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "494701 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1174,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5103553 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1332,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4497512 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1323,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4547098 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11903,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "494493 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1169,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5133655 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1341,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4465612 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1323,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4533684 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11120,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "534619 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1173,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5145001 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1332,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4510744 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1319,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4558639 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1183,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5076063 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1333,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4493732 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1325,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4540806 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3081,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1951086 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3115,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1875806 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3502,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1698537 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3173,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1895857 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5171,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Skye Gill",
            "username": "skyerus",
            "email": "gill.skye95@gmail.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "3f406b53bda8b5beb8b0929da3802a0368c13151",
          "message": "fix: omitempty targeting field in Flag structure (#247)\n\n## This PR\r\n<!-- add the description of the PR here -->\r\n\r\nAdded omitempty to the json tags of the `targeting` field in the `Flag`\r\nstructure.\r\n\r\n### Related Issues\r\n<!-- add here the GitHub issue that this PR resolves if applicable -->\r\n\r\nFixes #246 \r\n\r\n### Notes\r\n<!-- any additional notes for this PR -->\r\n\r\n### Follow-up Tasks\r\n<!-- anything that is related to this PR but not done here should be\r\nnoted under this section -->\r\n<!-- if there is a need for a new issue, please link it here -->\r\n\r\n### How to test\r\n<!-- if applicable, add testing instructions under this section -->\r\n\r\nSigned-off-by: Skye Gill <gill.skye95@gmail.com>",
          "timestamp": "2023-01-05T17:04:45Z",
          "url": "https://github.com/open-feature/flagd/commit/3f406b53bda8b5beb8b0929da3802a0368c13151"
        },
        "date": 1672972148669,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11944,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "492578 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1152,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5162617 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1342,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4470399 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1335,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4491115 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11890,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "500930 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1150,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5234900 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1331,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4484038 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1329,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4523914 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11977,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "491564 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1167,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5050341 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1325,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4524865 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1322,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4541690 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11016,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "536102 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1149,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5203407 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1327,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4521948 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1320,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4538002 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1179,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5068719 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1328,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4489834 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1335,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4490647 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3079,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1877133 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3107,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1932840 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3499,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1711819 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3164,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1897273 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5148,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "github-actions[bot]",
            "username": "github-actions[bot]",
            "email": "41898282+github-actions[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "2ad206e635cc5b68ae42287b740ddfc86b0194ff",
          "message": "chore(main): release 0.3.0 (#230)\n\n:robot: I have created a release *beep* *boop*\r\n---\r\n\r\n\r\n##\r\n[0.3.0](https://github.com/open-feature/flagd/compare/v0.2.7...v0.3.0)\r\n(2023-01-06)\r\n\r\n\r\n### ⚠ BREAKING CHANGES\r\n\r\n* consolidated configuration change events into one event\r\n([#241](https://github.com/open-feature/flagd/issues/241))\r\n\r\n### Features\r\n\r\n* consolidated configuration change events into one event\r\n([#241](https://github.com/open-feature/flagd/issues/241))\r\n([f9684b8](https://github.com/open-feature/flagd/commit/f9684b858dfef40576e0031654b421a37e8bb1d6))\r\n* support yaml evaluator\r\n([#206](https://github.com/open-feature/flagd/issues/206))\r\n([2dbace5](https://github.com/open-feature/flagd/commit/2dbace5b6bb8e187a7d44a3d3ec14190c63b3ae0))\r\n\r\n\r\n### Bug Fixes\r\n\r\n* changed eventing configuration mutex to rwmutex and added missing lock\r\n([#220](https://github.com/open-feature/flagd/issues/220))\r\n([5bbef9e](https://github.com/open-feature/flagd/commit/5bbef9ea4b1960686e58298c2c2e192ca99f072f))\r\n* omitempty targeting field in Flag structure\r\n([#247](https://github.com/open-feature/flagd/issues/247))\r\n([3f406b5](https://github.com/open-feature/flagd/commit/3f406b53bda8b5beb8b0929da3802a0368c13151))\r\n\r\n---\r\nThis PR was generated with [Release\r\nPlease](https://github.com/googleapis/release-please). See\r\n[documentation](https://github.com/googleapis/release-please#release-please).\r\n\r\nCo-authored-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2023-01-06T15:18:55Z",
          "url": "https://github.com/open-feature/flagd/commit/2ad206e635cc5b68ae42287b740ddfc86b0194ff"
        },
        "date": 1673058320684,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11695,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "497008 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1146,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5251315 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1288,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4675334 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1301,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4603044 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11590,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "510223 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1135,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5274867 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1293,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4638687 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1291,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4656124 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11758,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "506810 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1137,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5257700 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1302,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4602216 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1289,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4671212 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 10734,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "550368 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1138,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5277918 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1287,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4658134 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1293,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4634894 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1150,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5218832 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1295,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4644342 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1295,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4626631 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2893,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2089983 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2929,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "2052888 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3284,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1836416 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2958,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2023352 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4757,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1256119 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "github-actions[bot]",
            "username": "github-actions[bot]",
            "email": "41898282+github-actions[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "2ad206e635cc5b68ae42287b740ddfc86b0194ff",
          "message": "chore(main): release 0.3.0 (#230)\n\n:robot: I have created a release *beep* *boop*\r\n---\r\n\r\n\r\n##\r\n[0.3.0](https://github.com/open-feature/flagd/compare/v0.2.7...v0.3.0)\r\n(2023-01-06)\r\n\r\n\r\n### ⚠ BREAKING CHANGES\r\n\r\n* consolidated configuration change events into one event\r\n([#241](https://github.com/open-feature/flagd/issues/241))\r\n\r\n### Features\r\n\r\n* consolidated configuration change events into one event\r\n([#241](https://github.com/open-feature/flagd/issues/241))\r\n([f9684b8](https://github.com/open-feature/flagd/commit/f9684b858dfef40576e0031654b421a37e8bb1d6))\r\n* support yaml evaluator\r\n([#206](https://github.com/open-feature/flagd/issues/206))\r\n([2dbace5](https://github.com/open-feature/flagd/commit/2dbace5b6bb8e187a7d44a3d3ec14190c63b3ae0))\r\n\r\n\r\n### Bug Fixes\r\n\r\n* changed eventing configuration mutex to rwmutex and added missing lock\r\n([#220](https://github.com/open-feature/flagd/issues/220))\r\n([5bbef9e](https://github.com/open-feature/flagd/commit/5bbef9ea4b1960686e58298c2c2e192ca99f072f))\r\n* omitempty targeting field in Flag structure\r\n([#247](https://github.com/open-feature/flagd/issues/247))\r\n([3f406b5](https://github.com/open-feature/flagd/commit/3f406b53bda8b5beb8b0929da3802a0368c13151))\r\n\r\n---\r\nThis PR was generated with [Release\r\nPlease](https://github.com/googleapis/release-please). See\r\n[documentation](https://github.com/googleapis/release-please#release-please).\r\n\r\nCo-authored-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2023-01-06T15:18:55Z",
          "url": "https://github.com/open-feature/flagd/commit/2ad206e635cc5b68ae42287b740ddfc86b0194ff"
        },
        "date": 1673144961567,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 12024,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "493928 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1157,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5199483 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1327,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4509799 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1335,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4491379 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11916,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "495874 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1165,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5154278 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1339,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4482249 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1325,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4525028 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11994,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "455523 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1154,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5184638 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1324,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4526164 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1320,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4546004 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11078,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "534302 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1146,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5234036 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1327,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4525308 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1327,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4522794 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1167,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5135428 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1331,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4471869 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1321,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4540381 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3067,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1915290 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3111,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1926573 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3513,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1688047 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3195,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1883493 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5086,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "github-actions[bot]",
            "username": "github-actions[bot]",
            "email": "41898282+github-actions[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "2ad206e635cc5b68ae42287b740ddfc86b0194ff",
          "message": "chore(main): release 0.3.0 (#230)\n\n:robot: I have created a release *beep* *boop*\r\n---\r\n\r\n\r\n##\r\n[0.3.0](https://github.com/open-feature/flagd/compare/v0.2.7...v0.3.0)\r\n(2023-01-06)\r\n\r\n\r\n### ⚠ BREAKING CHANGES\r\n\r\n* consolidated configuration change events into one event\r\n([#241](https://github.com/open-feature/flagd/issues/241))\r\n\r\n### Features\r\n\r\n* consolidated configuration change events into one event\r\n([#241](https://github.com/open-feature/flagd/issues/241))\r\n([f9684b8](https://github.com/open-feature/flagd/commit/f9684b858dfef40576e0031654b421a37e8bb1d6))\r\n* support yaml evaluator\r\n([#206](https://github.com/open-feature/flagd/issues/206))\r\n([2dbace5](https://github.com/open-feature/flagd/commit/2dbace5b6bb8e187a7d44a3d3ec14190c63b3ae0))\r\n\r\n\r\n### Bug Fixes\r\n\r\n* changed eventing configuration mutex to rwmutex and added missing lock\r\n([#220](https://github.com/open-feature/flagd/issues/220))\r\n([5bbef9e](https://github.com/open-feature/flagd/commit/5bbef9ea4b1960686e58298c2c2e192ca99f072f))\r\n* omitempty targeting field in Flag structure\r\n([#247](https://github.com/open-feature/flagd/issues/247))\r\n([3f406b5](https://github.com/open-feature/flagd/commit/3f406b53bda8b5beb8b0929da3802a0368c13151))\r\n\r\n---\r\nThis PR was generated with [Release\r\nPlease](https://github.com/googleapis/release-please). See\r\n[documentation](https://github.com/googleapis/release-please#release-please).\r\n\r\nCo-authored-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2023-01-06T15:18:55Z",
          "url": "https://github.com/open-feature/flagd/commit/2ad206e635cc5b68ae42287b740ddfc86b0194ff"
        },
        "date": 1673231159983,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 12204,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "493437 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1190,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5007104 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1370,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4421756 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1340,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4491057 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 12059,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "489841 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1155,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5222916 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1342,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4471934 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1339,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4490064 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 12111,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "487268 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1166,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5119725 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1336,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4522393 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1318,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4510478 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11163,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "536176 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1163,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5141686 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1335,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4442539 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1336,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4482453 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1181,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5058536 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1347,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4468964 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1345,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4465926 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3128,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1953662 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3141,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1896139 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3530,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1689733 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3164,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1889012 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5158,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "github-actions[bot]",
            "username": "github-actions[bot]",
            "email": "41898282+github-actions[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "2ad206e635cc5b68ae42287b740ddfc86b0194ff",
          "message": "chore(main): release 0.3.0 (#230)\n\n:robot: I have created a release *beep* *boop*\r\n---\r\n\r\n\r\n##\r\n[0.3.0](https://github.com/open-feature/flagd/compare/v0.2.7...v0.3.0)\r\n(2023-01-06)\r\n\r\n\r\n### ⚠ BREAKING CHANGES\r\n\r\n* consolidated configuration change events into one event\r\n([#241](https://github.com/open-feature/flagd/issues/241))\r\n\r\n### Features\r\n\r\n* consolidated configuration change events into one event\r\n([#241](https://github.com/open-feature/flagd/issues/241))\r\n([f9684b8](https://github.com/open-feature/flagd/commit/f9684b858dfef40576e0031654b421a37e8bb1d6))\r\n* support yaml evaluator\r\n([#206](https://github.com/open-feature/flagd/issues/206))\r\n([2dbace5](https://github.com/open-feature/flagd/commit/2dbace5b6bb8e187a7d44a3d3ec14190c63b3ae0))\r\n\r\n\r\n### Bug Fixes\r\n\r\n* changed eventing configuration mutex to rwmutex and added missing lock\r\n([#220](https://github.com/open-feature/flagd/issues/220))\r\n([5bbef9e](https://github.com/open-feature/flagd/commit/5bbef9ea4b1960686e58298c2c2e192ca99f072f))\r\n* omitempty targeting field in Flag structure\r\n([#247](https://github.com/open-feature/flagd/issues/247))\r\n([3f406b5](https://github.com/open-feature/flagd/commit/3f406b53bda8b5beb8b0929da3802a0368c13151))\r\n\r\n---\r\nThis PR was generated with [Release\r\nPlease](https://github.com/googleapis/release-please). See\r\n[documentation](https://github.com/googleapis/release-please#release-please).\r\n\r\nCo-authored-by: github-actions[bot] <41898282+github-actions[bot]@users.noreply.github.com>",
          "timestamp": "2023-01-06T15:18:55Z",
          "url": "https://github.com/open-feature/flagd/commit/2ad206e635cc5b68ae42287b740ddfc86b0194ff"
        },
        "date": 1673317804635,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11642,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "511327 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1141,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5253274 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1296,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4629459 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1303,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4599632 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11538,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "510979 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1153,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5205183 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1306,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4587482 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1304,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4599393 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11613,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "508998 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1149,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5236032 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1307,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4585096 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1305,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4611154 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 10662,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "551862 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1155,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5187932 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1302,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4569326 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1302,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4624725 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1151,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5214222 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1306,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4586216 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1302,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4615184 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2891,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2036449 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2950,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "2032888 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3288,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1798710 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2969,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2024552 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4746,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1261774 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "renovate[bot]",
            "username": "renovate[bot]",
            "email": "29139614+renovate[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "5ddc8f79584206ae481083901521ed371dd764b2",
          "message": "chore(deps): update docker/login-action digest to 3da7dc6 (#254)\n\n[![Mend\nRenovate](https://app.renovatebot.com/images/banner.svg)](https://renovatebot.com)\n\nThis PR contains the following updates:\n\n| Package | Type | Update | Change |\n|---|---|---|---|\n| docker/login-action | action | digest | `f054a8b` -> `3da7dc6` |\n\n---\n\n### Configuration\n\n📅 **Schedule**: Branch creation - At any time (no schedule defined),\nAutomerge - At any time (no schedule defined).\n\n🚦 **Automerge**: Enabled.\n\n♻ **Rebasing**: Whenever PR becomes conflicted, or you tick the\nrebase/retry checkbox.\n\n🔕 **Ignore**: Close this PR and you won't be reminded about this update\nagain.\n\n---\n\n- [ ] <!-- rebase-check -->If you want to rebase/retry this PR, check\nthis box\n\n---\n\nThis PR has been generated by [Mend\nRenovate](https://www.mend.io/free-developer-tools/renovate/). View\nrepository job log\n[here](https://app.renovatebot.com/dashboard#github/open-feature/flagd).\n\n<!--renovate-debug:eyJjcmVhdGVkSW5WZXIiOiIzNC45Ny4xIiwidXBkYXRlZEluVmVyIjoiMzQuOTcuMSJ9-->\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2023-01-11T00:39:31Z",
          "url": "https://github.com/open-feature/flagd/commit/5ddc8f79584206ae481083901521ed371dd764b2"
        },
        "date": 1673404070817,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 12250,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "483121 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1163,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5194906 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1341,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4460088 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1334,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4485266 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 12051,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "492184 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1174,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5105544 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1357,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4424252 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1342,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4444527 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 12089,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "480813 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1157,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5188699 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1338,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4485156 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1324,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4517985 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11119,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "533754 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1151,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5213265 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1340,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4474198 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1334,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4251760 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1163,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5159145 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1342,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4458112 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1335,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4479068 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3113,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1947492 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3127,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1918892 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3534,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1701741 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3205,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1859316 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5142,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "renovate[bot]",
            "username": "renovate[bot]",
            "email": "29139614+renovate[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "2ec61c6bc61c266451b496ff18c3dd9a74173233",
          "message": "fix(deps): update module github.com/stretchr/testify to v1.8.1 (#265)\n\n[![Mend\nRenovate](https://app.renovatebot.com/images/banner.svg)](https://renovatebot.com)\n\nThis PR contains the following updates:\n\n| Package | Type | Update | Change |\n|---|---|---|---|\n| [github.com/stretchr/testify](https://togithub.com/stretchr/testify) |\nrequire | patch | `v1.8.0` -> `v1.8.1` |\n\n---\n\n### Release Notes\n\n<details>\n<summary>stretchr/testify</summary>\n\n###\n[`v1.8.1`](https://togithub.com/stretchr/testify/compare/v1.8.0...v1.8.1)\n\n[Compare\nSource](https://togithub.com/stretchr/testify/compare/v1.8.0...v1.8.1)\n\n</details>\n\n---\n\n### Configuration\n\n📅 **Schedule**: Branch creation - At any time (no schedule defined),\nAutomerge - At any time (no schedule defined).\n\n🚦 **Automerge**: Enabled.\n\n♻ **Rebasing**: Whenever PR becomes conflicted, or you tick the\nrebase/retry checkbox.\n\n🔕 **Ignore**: Close this PR and you won't be reminded about this update\nagain.\n\n---\n\n- [ ] <!-- rebase-check -->If you want to rebase/retry this PR, check\nthis box\n\n---\n\nThis PR has been generated by [Mend\nRenovate](https://www.mend.io/free-developer-tools/renovate/). View\nrepository job log\n[here](https://app.renovatebot.com/dashboard#github/open-feature/flagd).\n\n<!--renovate-debug:eyJjcmVhdGVkSW5WZXIiOiIzNC45Ny4xIiwidXBkYXRlZEluVmVyIjoiMzQuOTcuMSJ9-->\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2023-01-12T00:21:25Z",
          "url": "https://github.com/open-feature/flagd/commit/2ec61c6bc61c266451b496ff18c3dd9a74173233"
        },
        "date": 1673490438151,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11608,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "495021 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1147,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5229601 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1314,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4577028 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1313,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4565502 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11801,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "495393 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1154,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5198768 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1314,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4613724 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1313,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4578843 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 12078,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "461323 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1165,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5152041 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1313,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4550233 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1314,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4561663 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11021,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "540478 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1171,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5117994 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1297,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4604986 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1302,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4640744 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1163,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5174787 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1319,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4565398 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1297,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4580190 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2948,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2041185 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3012,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "2012742 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3369,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1782723 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3040,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1965180 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4763,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1254807 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "renovate[bot]",
            "username": "renovate[bot]",
            "email": "29139614+renovate[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "2ab5a00fa6f19c7e0fe1a4e36649fae2633ac269",
          "message": "fix(deps): update module github.com/diegoholiveira/jsonlogic/v3 to v3.2.7 (#283)\n\n[![Mend\nRenovate](https://app.renovatebot.com/images/banner.svg)](https://renovatebot.com)\n\nThis PR contains the following updates:\n\n| Package | Type | Update | Change |\n|---|---|---|---|\n|\n[github.com/diegoholiveira/jsonlogic/v3](https://togithub.com/diegoholiveira/jsonlogic)\n| require | patch | `v3.2.6` -> `v3.2.7` |\n\n---\n\n### Release Notes\n\n<details>\n<summary>diegoholiveira/jsonlogic</summary>\n\n###\n[`v3.2.7`](https://togithub.com/diegoholiveira/jsonlogic/releases/tag/v3.2.7)\n\n[Compare\nSource](https://togithub.com/diegoholiveira/jsonlogic/compare/v3.2.6...v3.2.7)\n\n#### What's Changed\n\n- Fixed subtraction from zero by\n[@&#8203;Rymmugygr](https://togithub.com/Rymmugygr) in\n[https://github.com/diegoholiveira/jsonlogic/pull/62](https://togithub.com/diegoholiveira/jsonlogic/pull/62)\n- Add the toSliceOfNumbers to clear the code in `-` and `/` functions by\n[@&#8203;diegoholiveira](https://togithub.com/diegoholiveira) in\n[https://github.com/diegoholiveira/jsonlogic/pull/63](https://togithub.com/diegoholiveira/jsonlogic/pull/63)\n- Fix null comparisons by\n[@&#8203;reb-felipe](https://togithub.com/reb-felipe) in\n[https://github.com/diegoholiveira/jsonlogic/pull/64](https://togithub.com/diegoholiveira/jsonlogic/pull/64)\n\n#### New Contributors\n\n- [@&#8203;Rymmugygr](https://togithub.com/Rymmugygr) made their first\ncontribution in\n[https://github.com/diegoholiveira/jsonlogic/pull/62](https://togithub.com/diegoholiveira/jsonlogic/pull/62)\n- [@&#8203;reb-felipe](https://togithub.com/reb-felipe) made their first\ncontribution in\n[https://github.com/diegoholiveira/jsonlogic/pull/64](https://togithub.com/diegoholiveira/jsonlogic/pull/64)\n\n**Full Changelog**:\nhttps://github.com/diegoholiveira/jsonlogic/compare/v3.2.6...v3.2.7\n\n</details>\n\n---\n\n### Configuration\n\n📅 **Schedule**: Branch creation - At any time (no schedule defined),\nAutomerge - At any time (no schedule defined).\n\n🚦 **Automerge**: Enabled.\n\n♻ **Rebasing**: Whenever PR becomes conflicted, or you tick the\nrebase/retry checkbox.\n\n🔕 **Ignore**: Close this PR and you won't be reminded about this update\nagain.\n\n---\n\n- [ ] <!-- rebase-check -->If you want to rebase/retry this PR, check\nthis box\n\n---\n\nThis PR has been generated by [Mend\nRenovate](https://www.mend.io/free-developer-tools/renovate/). View\nrepository job log\n[here](https://app.renovatebot.com/dashboard#github/open-feature/flagd).\n\n<!--renovate-debug:eyJjcmVhdGVkSW5WZXIiOiIzNC4xMDAuMSIsInVwZGF0ZWRJblZlciI6IjM0LjEwMC4xIn0=-->\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2023-01-13T00:59:40Z",
          "url": "https://github.com/open-feature/flagd/commit/2ab5a00fa6f19c7e0fe1a4e36649fae2633ac269"
        },
        "date": 1673577016504,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 12010,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "493051 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1155,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5200308 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1328,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4507226 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1321,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4543298 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 12005,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "489210 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1157,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5177010 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1333,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4503103 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1328,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4526487 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 12080,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "490791 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1168,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5135676 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1324,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4524139 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1326,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4523695 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11161,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "524612 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1170,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5119152 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1329,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4516404 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1324,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4515193 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1175,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5086210 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1326,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4529157 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1327,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4521702 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3231,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1875345 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3115,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1927274 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3547,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1694463 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3127,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1925857 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5209,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "renovate[bot]",
            "username": "renovate[bot]",
            "email": "29139614+renovate[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "f5b0eceffe6ff712b57e5bc2abbf584b05771107",
          "message": "chore(deps): update docker/metadata-action digest to 507c2f2 (#287)\n\n[![Mend\nRenovate](https://app.renovatebot.com/images/banner.svg)](https://renovatebot.com)\n\nThis PR contains the following updates:\n\n| Package | Type | Update | Change |\n|---|---|---|---|\n| docker/metadata-action | action | digest | `05d22bf` -> `507c2f2` |\n\n---\n\n### Configuration\n\n📅 **Schedule**: Branch creation - At any time (no schedule defined),\nAutomerge - At any time (no schedule defined).\n\n🚦 **Automerge**: Enabled.\n\n♻ **Rebasing**: Whenever PR becomes conflicted, or you tick the\nrebase/retry checkbox.\n\n🔕 **Ignore**: Close this PR and you won't be reminded about this update\nagain.\n\n---\n\n- [ ] <!-- rebase-check -->If you want to rebase/retry this PR, check\nthis box\n\n---\n\nThis PR has been generated by [Mend\nRenovate](https://www.mend.io/free-developer-tools/renovate/). View\nrepository job log\n[here](https://app.renovatebot.com/dashboard#github/open-feature/flagd).\n\n<!--renovate-debug:eyJjcmVhdGVkSW5WZXIiOiIzNC4xMDAuMSIsInVwZGF0ZWRJblZlciI6IjM0LjEwMC4xIn0=-->\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2023-01-13T23:24:30Z",
          "url": "https://github.com/open-feature/flagd/commit/f5b0eceffe6ff712b57e5bc2abbf584b05771107"
        },
        "date": 1673662931692,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 12083,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "491482 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1148,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5228548 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1333,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4490683 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1328,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4494300 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 12047,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "494025 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1154,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5146856 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1335,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4494976 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1330,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4524795 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 12102,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "487227 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1173,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5141978 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1332,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4495021 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1319,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4544208 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11210,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "527980 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1157,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5169393 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1331,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4503165 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1324,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4540030 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1162,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5161527 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1323,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4517949 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1337,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4496864 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3184,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1836298 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3124,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1917886 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3500,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1723668 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3108,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1929732 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5089,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "renovate[bot]",
            "username": "renovate[bot]",
            "email": "29139614+renovate[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "f5b0eceffe6ff712b57e5bc2abbf584b05771107",
          "message": "chore(deps): update docker/metadata-action digest to 507c2f2 (#287)\n\n[![Mend\nRenovate](https://app.renovatebot.com/images/banner.svg)](https://renovatebot.com)\n\nThis PR contains the following updates:\n\n| Package | Type | Update | Change |\n|---|---|---|---|\n| docker/metadata-action | action | digest | `05d22bf` -> `507c2f2` |\n\n---\n\n### Configuration\n\n📅 **Schedule**: Branch creation - At any time (no schedule defined),\nAutomerge - At any time (no schedule defined).\n\n🚦 **Automerge**: Enabled.\n\n♻ **Rebasing**: Whenever PR becomes conflicted, or you tick the\nrebase/retry checkbox.\n\n🔕 **Ignore**: Close this PR and you won't be reminded about this update\nagain.\n\n---\n\n- [ ] <!-- rebase-check -->If you want to rebase/retry this PR, check\nthis box\n\n---\n\nThis PR has been generated by [Mend\nRenovate](https://www.mend.io/free-developer-tools/renovate/). View\nrepository job log\n[here](https://app.renovatebot.com/dashboard#github/open-feature/flagd).\n\n<!--renovate-debug:eyJjcmVhdGVkSW5WZXIiOiIzNC4xMDAuMSIsInVwZGF0ZWRJblZlciI6IjM0LjEwMC4xIn0=-->\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2023-01-13T23:24:30Z",
          "url": "https://github.com/open-feature/flagd/commit/f5b0eceffe6ff712b57e5bc2abbf584b05771107"
        },
        "date": 1673749755562,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 13016,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "451605 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1322,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "4537819 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1534,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "3916561 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1490,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4017846 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 13296,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "466044 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1326,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4568828 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1563,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3834973 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1529,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3895120 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 13379,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "463970 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1313,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4568295 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1555,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3876889 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1491,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3965973 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 12251,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "486098 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1290,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4542794 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1558,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "3920952 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1500,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3960937 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1354,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4568446 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1511,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3845472 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1560,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3926203 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3291,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1794660 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3387,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1776762 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3777,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1626819 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3226,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1796474 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5389,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "renovate[bot]",
            "username": "renovate[bot]",
            "email": "29139614+renovate[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "f5b0eceffe6ff712b57e5bc2abbf584b05771107",
          "message": "chore(deps): update docker/metadata-action digest to 507c2f2 (#287)\n\n[![Mend\nRenovate](https://app.renovatebot.com/images/banner.svg)](https://renovatebot.com)\n\nThis PR contains the following updates:\n\n| Package | Type | Update | Change |\n|---|---|---|---|\n| docker/metadata-action | action | digest | `05d22bf` -> `507c2f2` |\n\n---\n\n### Configuration\n\n📅 **Schedule**: Branch creation - At any time (no schedule defined),\nAutomerge - At any time (no schedule defined).\n\n🚦 **Automerge**: Enabled.\n\n♻ **Rebasing**: Whenever PR becomes conflicted, or you tick the\nrebase/retry checkbox.\n\n🔕 **Ignore**: Close this PR and you won't be reminded about this update\nagain.\n\n---\n\n- [ ] <!-- rebase-check -->If you want to rebase/retry this PR, check\nthis box\n\n---\n\nThis PR has been generated by [Mend\nRenovate](https://www.mend.io/free-developer-tools/renovate/). View\nrepository job log\n[here](https://app.renovatebot.com/dashboard#github/open-feature/flagd).\n\n<!--renovate-debug:eyJjcmVhdGVkSW5WZXIiOiIzNC4xMDAuMSIsInVwZGF0ZWRJblZlciI6IjM0LjEwMC4xIn0=-->\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2023-01-13T23:24:30Z",
          "url": "https://github.com/open-feature/flagd/commit/f5b0eceffe6ff712b57e5bc2abbf584b05771107"
        },
        "date": 1673836081526,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11673,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "507925 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1133,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5314041 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1291,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4640457 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1290,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4647968 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11616,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "508461 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1146,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5235826 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1304,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4597118 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1283,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4618274 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11737,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "499300 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1141,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5260003 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1292,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4617156 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1290,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4639995 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 10705,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "551371 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1132,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5281212 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1290,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4637242 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1293,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4644646 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1139,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5252858 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1298,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4636297 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1297,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4616090 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3017,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2009562 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 2967,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "2020000 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3267,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1840186 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2938,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2043334 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5250,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1249126 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "renovate[bot]",
            "username": "renovate[bot]",
            "email": "29139614+renovate[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "f5b0eceffe6ff712b57e5bc2abbf584b05771107",
          "message": "chore(deps): update docker/metadata-action digest to 507c2f2 (#287)\n\n[![Mend\nRenovate](https://app.renovatebot.com/images/banner.svg)](https://renovatebot.com)\n\nThis PR contains the following updates:\n\n| Package | Type | Update | Change |\n|---|---|---|---|\n| docker/metadata-action | action | digest | `05d22bf` -> `507c2f2` |\n\n---\n\n### Configuration\n\n📅 **Schedule**: Branch creation - At any time (no schedule defined),\nAutomerge - At any time (no schedule defined).\n\n🚦 **Automerge**: Enabled.\n\n♻ **Rebasing**: Whenever PR becomes conflicted, or you tick the\nrebase/retry checkbox.\n\n🔕 **Ignore**: Close this PR and you won't be reminded about this update\nagain.\n\n---\n\n- [ ] <!-- rebase-check -->If you want to rebase/retry this PR, check\nthis box\n\n---\n\nThis PR has been generated by [Mend\nRenovate](https://www.mend.io/free-developer-tools/renovate/). View\nrepository job log\n[here](https://app.renovatebot.com/dashboard#github/open-feature/flagd).\n\n<!--renovate-debug:eyJjcmVhdGVkSW5WZXIiOiIzNC4xMDAuMSIsInVwZGF0ZWRJblZlciI6IjM0LjEwMC4xIn0=-->\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2023-01-13T23:24:30Z",
          "url": "https://github.com/open-feature/flagd/commit/f5b0eceffe6ff712b57e5bc2abbf584b05771107"
        },
        "date": 1673922407862,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 12101,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "492186 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1172,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5103056 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1340,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4478557 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1317,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4554896 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 12067,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "495165 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1168,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5117520 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1330,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4495416 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1330,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4521158 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 12151,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "484614 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1155,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5186523 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1338,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4462848 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1318,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4512867 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11247,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "522567 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1157,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5190885 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1337,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4495490 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1334,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4489152 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1158,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5182556 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1326,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4519322 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1323,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4524615 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3254,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1881999 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3132,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1919958 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3487,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1711962 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3110,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1921188 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5129,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "renovate[bot]",
            "username": "renovate[bot]",
            "email": "29139614+renovate[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "d042e1f4a3236a908b90b796df5a7095c58a98b2",
          "message": "chore(deps): update docker/metadata-action digest to 6c3ca5d (#289)\n\n[![Mend\nRenovate](https://app.renovatebot.com/images/banner.svg)](https://renovatebot.com)\n\nThis PR contains the following updates:\n\n| Package | Type | Update | Change |\n|---|---|---|---|\n| docker/metadata-action | action | digest | `507c2f2` -> `6c3ca5d` |\n\n---\n\n### Configuration\n\n📅 **Schedule**: Branch creation - At any time (no schedule defined),\nAutomerge - At any time (no schedule defined).\n\n🚦 **Automerge**: Enabled.\n\n♻ **Rebasing**: Whenever PR becomes conflicted, or you tick the\nrebase/retry checkbox.\n\n🔕 **Ignore**: Close this PR and you won't be reminded about this update\nagain.\n\n---\n\n- [ ] <!-- rebase-check -->If you want to rebase/retry this PR, check\nthis box\n\n---\n\nThis PR has been generated by [Mend\nRenovate](https://www.mend.io/free-developer-tools/renovate/). View\nrepository job log\n[here](https://app.renovatebot.com/dashboard#github/open-feature/flagd).\n\n<!--renovate-debug:eyJjcmVhdGVkSW5WZXIiOiIzNC4xMDUuMSIsInVwZGF0ZWRJblZlciI6IjM0LjEwNS4xIn0=-->\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2023-01-18T02:27:24Z",
          "url": "https://github.com/open-feature/flagd/commit/d042e1f4a3236a908b90b796df5a7095c58a98b2"
        },
        "date": 1674008969744,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 12079,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "494643 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1153,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5204518 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1323,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4531675 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1328,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4549504 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11958,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "495086 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1160,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5159300 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1331,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4515912 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1320,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4558544 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 12093,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "489600 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1159,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5178036 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1325,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4541386 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1322,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4524073 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11227,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "524259 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1166,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5154895 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1331,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4500481 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1319,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4545277 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1155,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5197047 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1330,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4505635 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1320,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4536960 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3243,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1876885 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3154,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1917202 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3518,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1702977 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3128,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1914990 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5130,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "renovate[bot]",
            "username": "renovate[bot]",
            "email": "29139614+renovate[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "9db4186482483d790169a887093f0419a01f6fdc",
          "message": "chore(deps): update github/codeql-action digest to a34ca99 (#292)\n\n[![Mend\nRenovate](https://app.renovatebot.com/images/banner.svg)](https://renovatebot.com)\n\nThis PR contains the following updates:\n\n| Package | Type | Update | Change |\n|---|---|---|---|\n| [github/codeql-action](https://togithub.com/github/codeql-action) |\naction | digest | `515828d` -> `a34ca99` |\n\n---\n\n### Configuration\n\n📅 **Schedule**: Branch creation - At any time (no schedule defined),\nAutomerge - At any time (no schedule defined).\n\n🚦 **Automerge**: Enabled.\n\n♻ **Rebasing**: Whenever PR becomes conflicted, or you tick the\nrebase/retry checkbox.\n\n🔕 **Ignore**: Close this PR and you won't be reminded about this update\nagain.\n\n---\n\n- [ ] <!-- rebase-check -->If you want to rebase/retry this PR, check\nthis box\n\n---\n\nThis PR has been generated by [Mend\nRenovate](https://www.mend.io/free-developer-tools/renovate/). View\nrepository job log\n[here](https://app.renovatebot.com/dashboard#github/open-feature/flagd).\n\n<!--renovate-debug:eyJjcmVhdGVkSW5WZXIiOiIzNC4xMDUuNCIsInVwZGF0ZWRJblZlciI6IjM0LjEwNS40In0=-->\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2023-01-18T23:55:42Z",
          "url": "https://github.com/open-feature/flagd/commit/9db4186482483d790169a887093f0419a01f6fdc"
        },
        "date": 1674095427595,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 12293,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "476854 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1155,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5174042 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1339,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4466545 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1321,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4508175 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 12121,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "491094 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1163,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5163090 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1351,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4418815 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1330,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4485060 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 12281,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "487328 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1163,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5153917 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1331,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4475104 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1338,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4493144 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11270,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "516100 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1184,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5154933 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1335,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4504074 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1330,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4494601 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1167,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5136876 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1332,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4506926 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1327,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4540540 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3267,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1856133 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3116,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1932862 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3587,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1669648 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3253,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1846154 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5221,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "renovate[bot]",
            "username": "renovate[bot]",
            "email": "29139614+renovate[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "d43220b2be58e4bce05050c5d1b36788289ae7cc",
          "message": "fix(deps): update module github.com/spf13/viper to v1.15.0 (#296)\n\n[![Mend\nRenovate](https://app.renovatebot.com/images/banner.svg)](https://renovatebot.com)\n\nThis PR contains the following updates:\n\n| Package | Type | Update | Change |\n|---|---|---|---|\n| [github.com/spf13/viper](https://togithub.com/spf13/viper) | require |\nminor | `v1.14.0` -> `v1.15.0` |\n\n---\n\n### Release Notes\n\n<details>\n<summary>spf13/viper</summary>\n\n### [`v1.15.0`](https://togithub.com/spf13/viper/releases/tag/v1.15.0)\n\n[Compare\nSource](https://togithub.com/spf13/viper/compare/v1.14.0...v1.15.0)\n\n<!-- Release notes generated using configuration in .github/release.yml\nat v1.15.0 -->\n\n#### What's Changed\n\n##### Exciting New Features 🎉\n\n- feat: add multiple endpoints support for remote by\n[@&#8203;mozartz](https://togithub.com/mozartz) in\n[https://github.com/spf13/viper/pull/1464](https://togithub.com/spf13/viper/pull/1464)\n\n##### Enhancements 🚀\n\n- Add DocBlock to WatchConfig by\n[@&#8203;glebik000](https://togithub.com/glebik000) in\n[https://github.com/spf13/viper/pull/1467](https://togithub.com/spf13/viper/pull/1467)\n\n##### Breaking Changes 🛠\n\n- Drop YAML v2 and TOML v1 by\n[@&#8203;sagikazarmark](https://togithub.com/sagikazarmark) in\n[https://github.com/spf13/viper/pull/1493](https://togithub.com/spf13/viper/pull/1493)\n- Drop support for Go 1.16 by\n[@&#8203;sagikazarmark](https://togithub.com/sagikazarmark) in\n[https://github.com/spf13/viper/pull/1494](https://togithub.com/spf13/viper/pull/1494)\n\n##### Dependency Updates ⬆️\n\n- build(deps): bump github.com/spf13/afero from 1.9.2 to 1.9.3 by\n[@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1465](https://togithub.com/spf13/viper/pull/1465)\n- build(deps): bump github.com/magiconair/properties from 1.8.6 to 1.8.7\nby [@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1475](https://togithub.com/spf13/viper/pull/1475)\n- build(deps): bump github.com/pelletier/go-toml/v2 from 2.0.5 to 2.0.6\nby [@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1466](https://togithub.com/spf13/viper/pull/1466)\n- build(deps): bump mheap/github-action-required-labels from 2 to 3 by\n[@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1482](https://togithub.com/spf13/viper/pull/1482)\n- build(deps): bump github.com/subosito/gotenv from 1.4.1 to 1.4.2 by\n[@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1488](https://togithub.com/spf13/viper/pull/1488)\n- build(deps): bump github.com/sagikazarmark/crypt from 0.8.0 to 0.9.0\nby [@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1490](https://togithub.com/spf13/viper/pull/1490)\n\n#### New Contributors\n\n- [@&#8203;choar816](https://togithub.com/choar816) made their first\ncontribution in\n[https://github.com/spf13/viper/pull/1481](https://togithub.com/spf13/viper/pull/1481)\n- [@&#8203;lol768](https://togithub.com/lol768) made their first\ncontribution in\n[https://github.com/spf13/viper/pull/1485](https://togithub.com/spf13/viper/pull/1485)\n- [@&#8203;mozartz](https://togithub.com/mozartz) made their first\ncontribution in\n[https://github.com/spf13/viper/pull/1464](https://togithub.com/spf13/viper/pull/1464)\n- [@&#8203;glebik000](https://togithub.com/glebik000) made their first\ncontribution in\n[https://github.com/spf13/viper/pull/1467](https://togithub.com/spf13/viper/pull/1467)\n\n**Full Changelog**:\nhttps://github.com/spf13/viper/compare/v1.14.0...v1.15.0\n\n</details>\n\n---\n\n### Configuration\n\n📅 **Schedule**: Branch creation - At any time (no schedule defined),\nAutomerge - At any time (no schedule defined).\n\n🚦 **Automerge**: Enabled.\n\n♻ **Rebasing**: Whenever PR becomes conflicted, or you tick the\nrebase/retry checkbox.\n\n🔕 **Ignore**: Close this PR and you won't be reminded about this update\nagain.\n\n---\n\n- [ ] <!-- rebase-check -->If you want to rebase/retry this PR, check\nthis box\n\n---\n\nThis PR has been generated by [Mend\nRenovate](https://www.mend.io/free-developer-tools/renovate/). View\nrepository job log\n[here](https://app.renovatebot.com/dashboard#github/open-feature/flagd).\n\n<!--renovate-debug:eyJjcmVhdGVkSW5WZXIiOiIzNC4xMDUuNCIsInVwZGF0ZWRJblZlciI6IjM0LjEwNS40In0=-->\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2023-01-19T23:38:52Z",
          "url": "https://github.com/open-feature/flagd/commit/d43220b2be58e4bce05050c5d1b36788289ae7cc"
        },
        "date": 1674181845569,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 13936,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "420033 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1361,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "4408915 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1584,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "3785763 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1574,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "3721158 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 13940,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "407530 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1393,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4249426 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1631,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3686589 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1629,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3685248 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 14443,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "362115 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1406,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4246311 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1604,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3707122 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1584,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3625886 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 12833,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "456711 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1389,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4346221 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1583,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "3771942 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1616,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3668870 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1422,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4172673 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1629,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3698041 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1600,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3722906 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3657,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1630858 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3572,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1670346 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3928,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1529533 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3574,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1690964 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5701,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "renovate[bot]",
            "username": "renovate[bot]",
            "email": "29139614+renovate[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "d43220b2be58e4bce05050c5d1b36788289ae7cc",
          "message": "fix(deps): update module github.com/spf13/viper to v1.15.0 (#296)\n\n[![Mend\nRenovate](https://app.renovatebot.com/images/banner.svg)](https://renovatebot.com)\n\nThis PR contains the following updates:\n\n| Package | Type | Update | Change |\n|---|---|---|---|\n| [github.com/spf13/viper](https://togithub.com/spf13/viper) | require |\nminor | `v1.14.0` -> `v1.15.0` |\n\n---\n\n### Release Notes\n\n<details>\n<summary>spf13/viper</summary>\n\n### [`v1.15.0`](https://togithub.com/spf13/viper/releases/tag/v1.15.0)\n\n[Compare\nSource](https://togithub.com/spf13/viper/compare/v1.14.0...v1.15.0)\n\n<!-- Release notes generated using configuration in .github/release.yml\nat v1.15.0 -->\n\n#### What's Changed\n\n##### Exciting New Features 🎉\n\n- feat: add multiple endpoints support for remote by\n[@&#8203;mozartz](https://togithub.com/mozartz) in\n[https://github.com/spf13/viper/pull/1464](https://togithub.com/spf13/viper/pull/1464)\n\n##### Enhancements 🚀\n\n- Add DocBlock to WatchConfig by\n[@&#8203;glebik000](https://togithub.com/glebik000) in\n[https://github.com/spf13/viper/pull/1467](https://togithub.com/spf13/viper/pull/1467)\n\n##### Breaking Changes 🛠\n\n- Drop YAML v2 and TOML v1 by\n[@&#8203;sagikazarmark](https://togithub.com/sagikazarmark) in\n[https://github.com/spf13/viper/pull/1493](https://togithub.com/spf13/viper/pull/1493)\n- Drop support for Go 1.16 by\n[@&#8203;sagikazarmark](https://togithub.com/sagikazarmark) in\n[https://github.com/spf13/viper/pull/1494](https://togithub.com/spf13/viper/pull/1494)\n\n##### Dependency Updates ⬆️\n\n- build(deps): bump github.com/spf13/afero from 1.9.2 to 1.9.3 by\n[@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1465](https://togithub.com/spf13/viper/pull/1465)\n- build(deps): bump github.com/magiconair/properties from 1.8.6 to 1.8.7\nby [@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1475](https://togithub.com/spf13/viper/pull/1475)\n- build(deps): bump github.com/pelletier/go-toml/v2 from 2.0.5 to 2.0.6\nby [@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1466](https://togithub.com/spf13/viper/pull/1466)\n- build(deps): bump mheap/github-action-required-labels from 2 to 3 by\n[@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1482](https://togithub.com/spf13/viper/pull/1482)\n- build(deps): bump github.com/subosito/gotenv from 1.4.1 to 1.4.2 by\n[@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1488](https://togithub.com/spf13/viper/pull/1488)\n- build(deps): bump github.com/sagikazarmark/crypt from 0.8.0 to 0.9.0\nby [@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1490](https://togithub.com/spf13/viper/pull/1490)\n\n#### New Contributors\n\n- [@&#8203;choar816](https://togithub.com/choar816) made their first\ncontribution in\n[https://github.com/spf13/viper/pull/1481](https://togithub.com/spf13/viper/pull/1481)\n- [@&#8203;lol768](https://togithub.com/lol768) made their first\ncontribution in\n[https://github.com/spf13/viper/pull/1485](https://togithub.com/spf13/viper/pull/1485)\n- [@&#8203;mozartz](https://togithub.com/mozartz) made their first\ncontribution in\n[https://github.com/spf13/viper/pull/1464](https://togithub.com/spf13/viper/pull/1464)\n- [@&#8203;glebik000](https://togithub.com/glebik000) made their first\ncontribution in\n[https://github.com/spf13/viper/pull/1467](https://togithub.com/spf13/viper/pull/1467)\n\n**Full Changelog**:\nhttps://github.com/spf13/viper/compare/v1.14.0...v1.15.0\n\n</details>\n\n---\n\n### Configuration\n\n📅 **Schedule**: Branch creation - At any time (no schedule defined),\nAutomerge - At any time (no schedule defined).\n\n🚦 **Automerge**: Enabled.\n\n♻ **Rebasing**: Whenever PR becomes conflicted, or you tick the\nrebase/retry checkbox.\n\n🔕 **Ignore**: Close this PR and you won't be reminded about this update\nagain.\n\n---\n\n- [ ] <!-- rebase-check -->If you want to rebase/retry this PR, check\nthis box\n\n---\n\nThis PR has been generated by [Mend\nRenovate](https://www.mend.io/free-developer-tools/renovate/). View\nrepository job log\n[here](https://app.renovatebot.com/dashboard#github/open-feature/flagd).\n\n<!--renovate-debug:eyJjcmVhdGVkSW5WZXIiOiIzNC4xMDUuNCIsInVwZGF0ZWRJblZlciI6IjM0LjEwNS40In0=-->\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2023-01-19T23:38:52Z",
          "url": "https://github.com/open-feature/flagd/commit/d43220b2be58e4bce05050c5d1b36788289ae7cc"
        },
        "date": 1674267939257,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11768,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "500979 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1150,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5193028 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1300,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4616451 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1300,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4624912 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11686,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "482031 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1138,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5272563 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1302,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4598136 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1305,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4598793 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11819,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "502770 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1141,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5249791 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1305,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4606590 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1287,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4666825 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 10786,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "546073 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1134,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5285881 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1291,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4648748 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1307,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4633292 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1148,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5224027 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1303,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4610394 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1290,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4637712 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3067,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1996873 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3045,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1993532 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3317,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1810582 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2982,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1988697 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4833,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1228743 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "renovate[bot]",
            "username": "renovate[bot]",
            "email": "29139614+renovate[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "d43220b2be58e4bce05050c5d1b36788289ae7cc",
          "message": "fix(deps): update module github.com/spf13/viper to v1.15.0 (#296)\n\n[![Mend\nRenovate](https://app.renovatebot.com/images/banner.svg)](https://renovatebot.com)\n\nThis PR contains the following updates:\n\n| Package | Type | Update | Change |\n|---|---|---|---|\n| [github.com/spf13/viper](https://togithub.com/spf13/viper) | require |\nminor | `v1.14.0` -> `v1.15.0` |\n\n---\n\n### Release Notes\n\n<details>\n<summary>spf13/viper</summary>\n\n### [`v1.15.0`](https://togithub.com/spf13/viper/releases/tag/v1.15.0)\n\n[Compare\nSource](https://togithub.com/spf13/viper/compare/v1.14.0...v1.15.0)\n\n<!-- Release notes generated using configuration in .github/release.yml\nat v1.15.0 -->\n\n#### What's Changed\n\n##### Exciting New Features 🎉\n\n- feat: add multiple endpoints support for remote by\n[@&#8203;mozartz](https://togithub.com/mozartz) in\n[https://github.com/spf13/viper/pull/1464](https://togithub.com/spf13/viper/pull/1464)\n\n##### Enhancements 🚀\n\n- Add DocBlock to WatchConfig by\n[@&#8203;glebik000](https://togithub.com/glebik000) in\n[https://github.com/spf13/viper/pull/1467](https://togithub.com/spf13/viper/pull/1467)\n\n##### Breaking Changes 🛠\n\n- Drop YAML v2 and TOML v1 by\n[@&#8203;sagikazarmark](https://togithub.com/sagikazarmark) in\n[https://github.com/spf13/viper/pull/1493](https://togithub.com/spf13/viper/pull/1493)\n- Drop support for Go 1.16 by\n[@&#8203;sagikazarmark](https://togithub.com/sagikazarmark) in\n[https://github.com/spf13/viper/pull/1494](https://togithub.com/spf13/viper/pull/1494)\n\n##### Dependency Updates ⬆️\n\n- build(deps): bump github.com/spf13/afero from 1.9.2 to 1.9.3 by\n[@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1465](https://togithub.com/spf13/viper/pull/1465)\n- build(deps): bump github.com/magiconair/properties from 1.8.6 to 1.8.7\nby [@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1475](https://togithub.com/spf13/viper/pull/1475)\n- build(deps): bump github.com/pelletier/go-toml/v2 from 2.0.5 to 2.0.6\nby [@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1466](https://togithub.com/spf13/viper/pull/1466)\n- build(deps): bump mheap/github-action-required-labels from 2 to 3 by\n[@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1482](https://togithub.com/spf13/viper/pull/1482)\n- build(deps): bump github.com/subosito/gotenv from 1.4.1 to 1.4.2 by\n[@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1488](https://togithub.com/spf13/viper/pull/1488)\n- build(deps): bump github.com/sagikazarmark/crypt from 0.8.0 to 0.9.0\nby [@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1490](https://togithub.com/spf13/viper/pull/1490)\n\n#### New Contributors\n\n- [@&#8203;choar816](https://togithub.com/choar816) made their first\ncontribution in\n[https://github.com/spf13/viper/pull/1481](https://togithub.com/spf13/viper/pull/1481)\n- [@&#8203;lol768](https://togithub.com/lol768) made their first\ncontribution in\n[https://github.com/spf13/viper/pull/1485](https://togithub.com/spf13/viper/pull/1485)\n- [@&#8203;mozartz](https://togithub.com/mozartz) made their first\ncontribution in\n[https://github.com/spf13/viper/pull/1464](https://togithub.com/spf13/viper/pull/1464)\n- [@&#8203;glebik000](https://togithub.com/glebik000) made their first\ncontribution in\n[https://github.com/spf13/viper/pull/1467](https://togithub.com/spf13/viper/pull/1467)\n\n**Full Changelog**:\nhttps://github.com/spf13/viper/compare/v1.14.0...v1.15.0\n\n</details>\n\n---\n\n### Configuration\n\n📅 **Schedule**: Branch creation - At any time (no schedule defined),\nAutomerge - At any time (no schedule defined).\n\n🚦 **Automerge**: Enabled.\n\n♻ **Rebasing**: Whenever PR becomes conflicted, or you tick the\nrebase/retry checkbox.\n\n🔕 **Ignore**: Close this PR and you won't be reminded about this update\nagain.\n\n---\n\n- [ ] <!-- rebase-check -->If you want to rebase/retry this PR, check\nthis box\n\n---\n\nThis PR has been generated by [Mend\nRenovate](https://www.mend.io/free-developer-tools/renovate/). View\nrepository job log\n[here](https://app.renovatebot.com/dashboard#github/open-feature/flagd).\n\n<!--renovate-debug:eyJjcmVhdGVkSW5WZXIiOiIzNC4xMDUuNCIsInVwZGF0ZWRJblZlciI6IjM0LjEwNS40In0=-->\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2023-01-19T23:38:52Z",
          "url": "https://github.com/open-feature/flagd/commit/d43220b2be58e4bce05050c5d1b36788289ae7cc"
        },
        "date": 1674354557909,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11897,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "519381 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1126,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5270036 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1300,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4697457 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1279,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4611854 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11533,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "527804 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1143,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5222395 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1285,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4547773 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1296,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4598200 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11757,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "489448 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1108,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5286402 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1309,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4605921 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1285,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4549810 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 10733,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "547563 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1110,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5256022 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1295,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4374900 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1310,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4826788 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1140,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5245912 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1284,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4570132 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1315,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4617166 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 2990,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1997703 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3009,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "2115198 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3219,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1835155 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 2944,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2048020 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4837,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1262277 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "renovate[bot]",
            "username": "renovate[bot]",
            "email": "29139614+renovate[bot]@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "d43220b2be58e4bce05050c5d1b36788289ae7cc",
          "message": "fix(deps): update module github.com/spf13/viper to v1.15.0 (#296)\n\n[![Mend\nRenovate](https://app.renovatebot.com/images/banner.svg)](https://renovatebot.com)\n\nThis PR contains the following updates:\n\n| Package | Type | Update | Change |\n|---|---|---|---|\n| [github.com/spf13/viper](https://togithub.com/spf13/viper) | require |\nminor | `v1.14.0` -> `v1.15.0` |\n\n---\n\n### Release Notes\n\n<details>\n<summary>spf13/viper</summary>\n\n### [`v1.15.0`](https://togithub.com/spf13/viper/releases/tag/v1.15.0)\n\n[Compare\nSource](https://togithub.com/spf13/viper/compare/v1.14.0...v1.15.0)\n\n<!-- Release notes generated using configuration in .github/release.yml\nat v1.15.0 -->\n\n#### What's Changed\n\n##### Exciting New Features 🎉\n\n- feat: add multiple endpoints support for remote by\n[@&#8203;mozartz](https://togithub.com/mozartz) in\n[https://github.com/spf13/viper/pull/1464](https://togithub.com/spf13/viper/pull/1464)\n\n##### Enhancements 🚀\n\n- Add DocBlock to WatchConfig by\n[@&#8203;glebik000](https://togithub.com/glebik000) in\n[https://github.com/spf13/viper/pull/1467](https://togithub.com/spf13/viper/pull/1467)\n\n##### Breaking Changes 🛠\n\n- Drop YAML v2 and TOML v1 by\n[@&#8203;sagikazarmark](https://togithub.com/sagikazarmark) in\n[https://github.com/spf13/viper/pull/1493](https://togithub.com/spf13/viper/pull/1493)\n- Drop support for Go 1.16 by\n[@&#8203;sagikazarmark](https://togithub.com/sagikazarmark) in\n[https://github.com/spf13/viper/pull/1494](https://togithub.com/spf13/viper/pull/1494)\n\n##### Dependency Updates ⬆️\n\n- build(deps): bump github.com/spf13/afero from 1.9.2 to 1.9.3 by\n[@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1465](https://togithub.com/spf13/viper/pull/1465)\n- build(deps): bump github.com/magiconair/properties from 1.8.6 to 1.8.7\nby [@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1475](https://togithub.com/spf13/viper/pull/1475)\n- build(deps): bump github.com/pelletier/go-toml/v2 from 2.0.5 to 2.0.6\nby [@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1466](https://togithub.com/spf13/viper/pull/1466)\n- build(deps): bump mheap/github-action-required-labels from 2 to 3 by\n[@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1482](https://togithub.com/spf13/viper/pull/1482)\n- build(deps): bump github.com/subosito/gotenv from 1.4.1 to 1.4.2 by\n[@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1488](https://togithub.com/spf13/viper/pull/1488)\n- build(deps): bump github.com/sagikazarmark/crypt from 0.8.0 to 0.9.0\nby [@&#8203;dependabot](https://togithub.com/dependabot) in\n[https://github.com/spf13/viper/pull/1490](https://togithub.com/spf13/viper/pull/1490)\n\n#### New Contributors\n\n- [@&#8203;choar816](https://togithub.com/choar816) made their first\ncontribution in\n[https://github.com/spf13/viper/pull/1481](https://togithub.com/spf13/viper/pull/1481)\n- [@&#8203;lol768](https://togithub.com/lol768) made their first\ncontribution in\n[https://github.com/spf13/viper/pull/1485](https://togithub.com/spf13/viper/pull/1485)\n- [@&#8203;mozartz](https://togithub.com/mozartz) made their first\ncontribution in\n[https://github.com/spf13/viper/pull/1464](https://togithub.com/spf13/viper/pull/1464)\n- [@&#8203;glebik000](https://togithub.com/glebik000) made their first\ncontribution in\n[https://github.com/spf13/viper/pull/1467](https://togithub.com/spf13/viper/pull/1467)\n\n**Full Changelog**:\nhttps://github.com/spf13/viper/compare/v1.14.0...v1.15.0\n\n</details>\n\n---\n\n### Configuration\n\n📅 **Schedule**: Branch creation - At any time (no schedule defined),\nAutomerge - At any time (no schedule defined).\n\n🚦 **Automerge**: Enabled.\n\n♻ **Rebasing**: Whenever PR becomes conflicted, or you tick the\nrebase/retry checkbox.\n\n🔕 **Ignore**: Close this PR and you won't be reminded about this update\nagain.\n\n---\n\n- [ ] <!-- rebase-check -->If you want to rebase/retry this PR, check\nthis box\n\n---\n\nThis PR has been generated by [Mend\nRenovate](https://www.mend.io/free-developer-tools/renovate/). View\nrepository job log\n[here](https://app.renovatebot.com/dashboard#github/open-feature/flagd).\n\n<!--renovate-debug:eyJjcmVhdGVkSW5WZXIiOiIzNC4xMDUuNCIsInVwZGF0ZWRJblZlciI6IjM0LjEwNS40In0=-->\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2023-01-19T23:38:52Z",
          "url": "https://github.com/open-feature/flagd/commit/d43220b2be58e4bce05050c5d1b36788289ae7cc"
        },
        "date": 1674440781894,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 12122,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "492445 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1167,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5122168 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1327,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4523972 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1327,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4511467 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 12161,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "485456 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1163,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5148765 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1336,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4478760 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1326,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4528251 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 12292,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "470028 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1162,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5151457 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1329,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4527633 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1323,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4518673 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 11253,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "515184 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1163,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5159554 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1333,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4494979 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1333,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4486167 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1165,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5144238 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1341,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4493366 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1316,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4563672 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3211,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1827282 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3117,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1923673 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3488,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1720395 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3111,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1931800 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5160,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1000000 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Kavindu Dodanduwa",
            "username": "Kavindu-Dodan",
            "email": "Kavindu-Dodan@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "46c5bf87267de8120a8de39c9e35f2ceb6f6034f",
          "message": "refactor: refactor sync providers (#291)\n\n## This PR\r\n\r\nRevisits sync providers and refactor them.\r\n\r\n- unify event handling\r\n- packaging improvements\r\n- removal of unused codes (cleanup)\r\n- clean contracts between Runtime and Sync implementations \r\n\r\nWith this change, we get a clear isolation between runtime and sync\r\nproviders (ex:- file, k8s, ....). For example, the runtime is not aware\r\nof any sync implementation details such as events. It simply coordinates\r\nsync impls, store (evaluator) and notifications. This should bring more\r\nextensibility when adding future extensions.\r\n\r\nBelow is the overview of the internals and interactions,\r\n\r\n\r\n![image](https://user-images.githubusercontent.com/8186721/213037032-316adb7e-e9ab-42e3-82f6-c3eaa2612ba3.png)\r\n\r\nSigned-off-by: Kavindu Dodanduwa <kavindudodanduwa@gmail.com>",
          "timestamp": "2023-01-23T16:05:54Z",
          "url": "https://github.com/open-feature/flagd/commit/46c5bf87267de8120a8de39c9e35f2ceb6f6034f"
        },
        "date": 1674527187366,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_staticBoolFlag",
            "value": 1205,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "4962925 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 11623,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "503145 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1138,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "5278057 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1283,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4669725 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1288,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4621384 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticStringFlag",
            "value": 1256,
            "unit": "ns/op\t     128 B/op\t       6 allocs/op",
            "extra": "4775912 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 11627,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "494890 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1144,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5266500 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1304,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4577302 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1278,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4684755 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticFloatFlag",
            "value": 1246,
            "unit": "ns/op\t     112 B/op\t       6 allocs/op",
            "extra": "4819390 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 11658,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "500865 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1138,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5253006 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1285,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4663455 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1308,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4580979 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticIntFlag",
            "value": 1202,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5004009 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 10709,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "547464 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1127,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5307543 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1289,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "4645758 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1287,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4666526 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticObjectFlag",
            "value": 4570,
            "unit": "ns/op\t    1392 B/op\t      32 allocs/op",
            "extra": "1311788 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_targetingObjectFlag",
            "value": 15007,
            "unit": "ns/op\t    6106 B/op\t     104 allocs/op",
            "extra": "396169 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1142,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "5270992 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1291,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4615485 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1293,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "4628931 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3028,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1987843 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3021,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1990500 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3311,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1817278 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3005,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "2006989 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 4811,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "1239546 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "Michael Beemer",
            "username": "beeme1mr",
            "email": "beeme1mr@users.noreply.github.com"
          },
          "committer": {
            "name": "GitHub",
            "username": "web-flow",
            "email": "noreply@github.com"
          },
          "id": "f9bf974059a46f6e3f44bcb12064284c8a09cc1c",
          "message": "doc: update binary start command\n\nSigned-off-by: Michael Beemer <beeme1mr@users.noreply.github.com>",
          "timestamp": "2023-01-25T02:26:20Z",
          "url": "https://github.com/open-feature/flagd/commit/f9bf974059a46f6e3f44bcb12064284c8a09cc1c"
        },
        "date": 1674613670870,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkResolveBooleanValue/test_staticBoolFlag",
            "value": 1452,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "4081987 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_targetingBoolFlag",
            "value": 13952,
            "unit": "ns/op\t    4801 B/op\t      80 allocs/op",
            "extra": "397273 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_staticObjectFlag",
            "value": 1404,
            "unit": "ns/op\t      80 B/op\t       4 allocs/op",
            "extra": "3950646 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_missingFlag",
            "value": 1617,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "3762601 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveBooleanValue/test_disabledFlag",
            "value": 1581,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "3777428 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticStringFlag",
            "value": 1566,
            "unit": "ns/op\t     128 B/op\t       6 allocs/op",
            "extra": "3743389 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_targetingStringFlag",
            "value": 13942,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "405826 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_staticObjectFlag",
            "value": 1410,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4287018 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_missingFlag",
            "value": 1622,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3582302 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveStringValue/test_disabledFlag",
            "value": 1585,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3602950 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticFloatFlag",
            "value": 1491,
            "unit": "ns/op\t     112 B/op\t       6 allocs/op",
            "extra": "4192808 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_targetingFloatFlag",
            "value": 13650,
            "unit": "ns/op\t    4841 B/op\t      82 allocs/op",
            "extra": "457382 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_staticObjectFlag",
            "value": 1395,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4282635 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_missingFlag",
            "value": 1611,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3572803 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveFloatValue/test_disabledFlag",
            "value": 1605,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3667941 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticIntFlag",
            "value": 1552,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4100912 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_targetingNumberFlag",
            "value": 13475,
            "unit": "ns/op\t    4825 B/op\t      80 allocs/op",
            "extra": "459922 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_staticObjectFlag",
            "value": 1399,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4410501 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_missingFlag",
            "value": 1626,
            "unit": "ns/op\t     144 B/op\t       6 allocs/op",
            "extra": "3790941 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveIntValue/test_disabledFlag",
            "value": 1544,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3853941 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticObjectFlag",
            "value": 5394,
            "unit": "ns/op\t    1392 B/op\t      32 allocs/op",
            "extra": "1000000 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_targetingObjectFlag",
            "value": 17588,
            "unit": "ns/op\t    6106 B/op\t     104 allocs/op",
            "extra": "320368 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_staticBoolFlag",
            "value": 1337,
            "unit": "ns/op\t      96 B/op\t       4 allocs/op",
            "extra": "4502894 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_missingFlag",
            "value": 1577,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3789634 times\n2 procs"
          },
          {
            "name": "BenchmarkResolveObjectValue/test_disabledFlag",
            "value": 1557,
            "unit": "ns/op\t     160 B/op\t       6 allocs/op",
            "extra": "3838129 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveBoolean/happy_path",
            "value": 3507,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1761914 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveString/happy_path",
            "value": 3425,
            "unit": "ns/op\t     584 B/op\t      15 allocs/op",
            "extra": "1750810 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveFloat/happy_path",
            "value": 3877,
            "unit": "ns/op\t     624 B/op\t      15 allocs/op",
            "extra": "1515296 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveInt/happy_path",
            "value": 3727,
            "unit": "ns/op\t     552 B/op\t      14 allocs/op",
            "extra": "1694643 times\n2 procs"
          },
          {
            "name": "BenchmarkConnectService_ResolveObject/happy_path",
            "value": 5866,
            "unit": "ns/op\t    1856 B/op\t      34 allocs/op",
            "extra": "972308 times\n2 procs"
          }
        ]
      }
    ]
  }
}