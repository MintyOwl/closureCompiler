package closureCompiler

const codee = `"use strict";

var helpersNameSpace = helpersNameSpace || {};
//let allocated = [2,5,9,12]
//var allocated = [3, 7, 9, 12];
//allocated.sort(function(a, b){return b-a});

helpersNameSpace.getPoolID = function (current) {
  var allocated = [];
  var sortNumber = function sortNumber(a, b) {
    return a - b;
  };

  allocated = current;
  allocated.sort(sortNumber);
  var pool = [];
  //pool.sort(function(a, b){return a-b});
  var sum = 0;

  var calc = function calc(prev, context) {
    //console.log(prev, context, context - prev - 1);
    for (var i = 1; i <= context - prev - 1; i++) {
      //console.log("Called ", i, " many times", "WITH", context - i);
      pool.push(context - i);
    }
  };

  var run = function run() {
    for (var i = 0; i < allocated.length; i++) {
      // console.log("SUM",sum)
      if (allocated[i] != i + 1) {
        var prev = 0;
        if (i != 0) {
          prev = allocated[i - 1];
        }
        calc(prev, allocated[i]);
      }
      sum = sum + allocated[i];
    }
    return pool.sort(sortNumber);
  };
  return run();
};

var newPool = helpersNameSpace.getPoolID([15, 3, 7, 9, 12]);
console.log(newPool;`

const jsn = `{  
   "compiledCode":"",
   "errors":[  
      {  
         "type":"JSC_PARSE_ERROR",
         "file":"Input_0",
         "lineno":46,
         "charno":19,
         "error":"Parse error. ',' expected",
         "line":"console.log(newPool;"
      }
   ],
   "statistics":{  
      "originalSize":1160,
      "originalGzipSize":488,
      "compressedSize":0,
      "compressedGzipSize":20,
      "compileTime":0
   },
   "outputFilePath":"/code/jscc4b3b3b87eb311cd4086bba6955a10c2/default.js"
}`
