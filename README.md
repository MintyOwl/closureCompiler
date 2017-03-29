# closureCompiler
Minify JS code using closure-compiler.appspot.com

````
jsCode = getSomeRawJSCode()
<<<<<<< HEAD
ua := "Custom User Agent"
=======
ua = "Custom User Agent"
>>>>>>> aa3e5f75b813d699a18fde42c10a0a9cd57fc381
cce = closureCompiler.NewCCEval(jsCode, ua)
minified, err = cce.Run() // minified is the result from closureCompiler
fmt.Println(minified, err)
````