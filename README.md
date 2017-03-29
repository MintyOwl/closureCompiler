# closureCompiler
Minify JS code using closure-compiler.appspot.com

````
jsCode = getSomeRawJSCode()
cce = closureCompiler.NewCCEval(jsCode, *ua)
		minified, err = cce.Run() // minified is the result from closureCompiler
````