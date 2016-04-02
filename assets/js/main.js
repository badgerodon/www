var Main = {
  api: function(path, args) {
    var tries = 5;
    function run() {
      tries--;
      $.ajax({
        "url": "/api/" + path,
        "data": args.data,
        "success": function() {
          if (args.success) {
            args.success.apply(this, arguments);
          }
        },
        "error": function(xhr) {
          if (xhr.status === 0 && tries) {
            // retry
            window.setTimeout(run, 1000);
            return;
          }
          if (args.error) {
            args.error.apply(this, arguments);
          }
        },
        "dataType": "json"
      });
    }
    run();
  }
};
