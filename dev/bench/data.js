window.BENCHMARK_DATA = {
  "lastUpdate": 1663988265645,
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
      }
    ]
  }
}