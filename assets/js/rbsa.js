// Load the Visualization API and the piechart package.
google.load('visualization', '1.0', {'packages':['corechart']});



$(function() {
  function drawChart(result) {
    // Create the data table.
    var data = new google.visualization.DataTable();
    data.addColumn('string', 'Symbol');
    data.addColumn('number', 'Amount');

    var rows = [];
    for (var key in result) {
      rows.push([key, Math.abs(result[key])]);
    }

    data.addRows(rows);

    // Set chart options
    var options = {
      'backgroundColor': '#FFF',
      'legend': {
        'position': 'left',
        'textStyle': {
          'color': '#000'
        }
      },
      'width': 350,
      'height': 200
    };

    // Instantiate and draw our chart, passing in some options.
    var $chart = $("<div class='chart'>");
    var chart = new google.visualization.PieChart($chart[0]);
    chart.draw(data, options);
    $("#rbsa-form").append($chart);
  }

  function drawTable(result) {
    var $table = $("<table><thead><tr><th>Symbol</th><th>Amount</th></tr></thead><tbody></tbody></table>");
    for (var key in result.data) {
      $table.find("tbody").append([
        '<tr>',
          '<td>',
            result.indices[key], ' (', key, ')',
          '</td>',
          '<td class="amount">',
            result.data[key] < 0.00001
            ? '0.00'
            : (result.data[key]*100).toFixed(2),
          '</td>',
        '</tr>'
      ].join(''));
    }
    $("#rbsa-form").append($table);
  }

  var ps = (location.search || "?").substr(1).split("&");
  var args = {};
  for (var i=0; i<ps.length; i++) {
    var nv = ps[i].split("=");
    args[decodeURIComponent(nv[0])] = decodeURIComponent(nv[1]) || "on";
  }

  if (args.symbol) {
    $("#rbsa-form").append("<h2>Loading...</h2>");
    Main.api("/rbsa", {
      "data": {
        "symbol": args.symbol
      },
      "success": function(result) {
        $("#rbsa-form h2").remove();
        $("#rbsa-form").append("<h2>Analysis of <i>" + args.symbol.toUpperCase() + "</i></h2>");
        drawTable(result);
        drawChart(result.data);
      },
      "error": function(xhr, err) {
        $("#rbsa-form").append('<div class="error">' + xhr.responseText + '</div>');
      }
    })
  }
});
