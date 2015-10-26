package yaag

const Template = `<!DOCTYPE html>
<html>
<head lang="en">
    <title> API Documentation </title>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap.min.css">
    <script src="http://google-code-prettify.googlecode.com/svn/loader/run_prettify.js"></script>
    <link href='http://fonts.googleapis.com/css?family=Roboto' rel='stylesheet' type='text/css'>
    <link rel="stylesheet" href="http://google-code-prettify.googlecode.com/svn/trunk/src/prettify.css">
    <!-- Optional theme -->
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap-theme.min.css">
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>
    <!-- Latest compiled and minified JavaScript -->
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/js/bootstrap.min.js"></script>
    <style type="text/css">
        body {
            font-family: 'Roboto', sans-serif;
        }
        .hidden {
            display:none;
        }
    </style>
    <style type="text/css">
        pre.prettyprint {
            border: 1px solid #ccc;
            margin-bottom: 0;
            padding: 9.5px;
        }
    </style>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/8.4/highlight.min.js"></script>
    <script>
        hljs.initHighlightingOnLoad();
        function toggler(divId) {
           $("#" + divId).toggle();
        }
    </script>
</head>
<body>
<nav class="navbar navbar-default navbar-fixed-top">
    <div class="container-fluid">
        <!-- Brand and toggle get grouped for better mobile display -->
        <div class="navbar-header">
            <button type="button" class="navbar-toggle collapsed" data-toggle="collapse"
                    data-target="#bs-example-navbar-collapse-1">
                <span class="sr-only">Toggle navigation</span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
            </button>
            <a class="navbar-brand" href="#">{{.Title}}</a>
            <p class="navbar-text">Developed by Gophers at <a href="http://betacraft.co">Betacraft</a></p>
        </div>
            
        <!-- /.navbar-collapse -->
    </div>
    <!-- /.container-fluid -->
</nav>
<div class="container-fluid" style="margin-top: 70px;margin-bottom: 20px;">
    <div class="container-fluid">
    <div class="col-md-4">
        <div class="panel panel-default">
              <div class="panel-heading">Base Urls</div>
              <div class="panel-body">
                {{ range $key, $value := .baseUrls }}
                    <p>{{$key}} : <strong>{{ $value }}</strong></p>
                {{ end }}
              </div>
            </div>    
        <ul class="nav nav-pills nav-stacked" role="tablist">
            {{ range $key, $value := .array }}
                <li role="presentation"><a href="#{{$key}}top" role="tab" data-toggle="tab">{{$value.HttpVerb}} : {{$value.Path}}</a></li>
            {{ end }}
        <ul>
    </div>
    <div class="col-md-8 tab-content">
        {{ range $key, $value := .array }}
        <div id="{{$key}}top"  role="tabpanel" class="tab-pane col-md-10">
            {{ range $wrapperKey, $wrapperValue := $value.Calls }}
                {{ if $wrapperValue.RequestHeader }}
                <p> <H4> Request Headers </H4> </p>
                <table class="table table-bordered table-striped">
                    <tr>
                        <th>Key</th>
                        <th>Value</th>
                    </tr>
                    {{ range $key, $value := $wrapperValue.RequestHeader }}
                    <tr>
                        <td>{{ $key }}</td>
                        <td> {{ $value }}</td>
                    </tr>
                    {{ end }}
                </table>
                {{ end }}
                {{ if $wrapperValue.PostForm }}
                <p> <H4> Post Form </H4> </p>
                <table class="table table-bordered table-striped">
                    <tr>
                        <th>Key</th>
                        <th>Value</th>
                    </tr>
                    {{ range $key, $value := $wrapperValue.PostForm }}
                    <tr>
                        <td>{{ $key }}</td>
                        <td> {{ $value }}</td>
                    </tr>
                    {{ end }}
                </table>
                {{ end }}
                {{ if $wrapperValue.RequestUrlParams }}
                <p> <H4> URL Params </H4> </p>
                <table class="table table-bordered table-striped">
                    <tr>
                        <th>Key</th>
                        <th>Value</th>
                    </tr>
                    {{ range $key, $value := $wrapperValue.RequestUrlParams }}
                    <tr>
                        <td>{{ $key }}</td>
                        <td> {{ $value }}</td>
                    </tr>
                    {{ end }}
                </table>
                {{ end }}
                {{ if $wrapperValue.RequestBody }}
                <p> <H4> Request Body </H4> </p>                
                <pre class="prettyprint lang-json">{{ $wrapperValue.RequestBody }}</pre>
                {{ end }}
                <p><h4> Response Code</h4></p>
                <pre class="prettyprint lang-json">{{ $wrapperValue.ResponseCode }}</pre>
                {{ if $wrapperValue.ResponseHeader }}
                <p><h4> Response Headers</h4></p>
                <table class="table table-bordered table-striped">
                    <tr>
                        <th>Key</th>
                        <th>Value</th>
                    </tr>
                    {{ range $key, $value := $wrapperValue.ResponseHeader }}
                    <tr>
                        <td>{{ $key }}</td>
                        <td> {{ $value }}</td>
                    </tr>
                    {{ end }}
                </table>
                {{ end }}
                {{ if $wrapperValue.ResponseBody }}
                <p> <H4> Response Body </H4> </p>
                <pre class="prettyprint lang-json">{{ $wrapperValue.ResponseBody }}</pre>
                {{ end }}
                <hr>
            {{ end }}
        </div>
        {{ end }}
    </div>   
</div>
</div>
<hr>
</body>
</html>`
