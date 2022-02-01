/**
 * Theme: Metrica - Responsive Bootstrap 4 Admin Dashboard
 * Author: Mannatthemes
 * Dashboard Js
 */


    //colunm-1

var options = {
        chart: {
            height: 340,
            type: 'bar',
            toolbar: {
                show: false
            },
            dropShadow: {
                enabled: true,
                top: 0,
                left: 5,
                bottom: 5,
                right: 0,
                blur: 5,
                color: '#b6c2e4',
                opacity: 0.35
            },
        },
        plotOptions: {
            bar: {
                horizontal: false,
                endingShape: 'rounded',
                columnWidth: '25%',
            },
        },
        dataLabels: {
            enabled: false,
        },
        stroke: {
            show: true,
            width: 2,
            colors: ['transparent']
        },
        colors: ["#2c77f4", "#1ecab8"],
        series: [{
            name: 'New Visitors',
            data: [68, 44, 55, 57, 56, 61, 58, 63, 60, 66]
        }, {
            name: 'Unique Visitors',
            data: [51, 76, 85, 101, 98, 87, 105, 91, 114, 94]
        },],
        xaxis: {
            categories: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct'],
            axisBorder: {
                show: true,
                color: '#bec7e0',
            },
            axisTicks: {
                show: true,
                color: '#bec7e0',
            },
        },
        legend: {
            offsetY: -10,
        },
        yaxis: {
            title: {
                text: 'Visitors'
            }
        },
        fill: {
            opacity: 1,
        },
        // legend: {
        //     floating: true
        // },
        grid: {
            row: {
                colors: ['transparent', 'transparent'], // takes an array which will be repeated on columns
                opacity: 0.2
            },
            borderColor: '#f1f3fa'
        },
        tooltip: {
            y: {
                formatter: function (val) {
                    return "" + val + "k"
                }
            }
        }
    }

var chart = new ApexCharts(
    document.querySelector("#ana_dash_1"),
    options
);

chart.render();


// traffice chart


var optionsCircle = {
    chart: {
        type: 'radialBar',
        height: 240,
        offsetY: -30,
        offsetX: 20,
        dropShadow: {
            enabled: true,
            top: 5,
            left: 0,
            bottom: 0,
            right: 0,
            blur: 5,
            color: '#b6c2e4',
            opacity: 0.35
        },
    },
    plotOptions: {
        radialBar: {
            inverseOrder: true,
            hollow: {
                margin: 5,
                size: '55%',
                background: 'transparent',
            },
            track: {
                show: true,
                background: '#ddd',
                strokeWidth: '10%',
                opacity: 1,
                margin: 5, // margin is in pixels
            },

            dataLabels: {
                name: {
                    fontSize: '18px',
                },
                value: {
                    fontSize: '16px',
                    color: '#50649c',
                },
            }
        },
    },
    series: [0],
    labels: ['Usage'],
    legend: {
        show: true,
        position: 'bottom',
        offsetX: -40,
        offsetY: -10,
        formatter: function (val, opts) {
            return val + " - " + opts.w.globals.series[opts.seriesIndex] + '%'
        }
    },
    fill: {
        type: 'gradient',
        gradient: {
            shade: 'dark',
            type: 'horizontal',
            shadeIntensity: 0.5,
            inverseColors: true,
            opacityFrom: 1,
            opacityTo: 1,
            stops: [0, 100],
            gradientToColors: ["#ffb822", "#5d78ff", "#34bfa3"]
        }
    },
    stroke: {
        lineCap: 'round'
    },
}
var optionsCircle2 = {
    chart: {
        type: 'radialBar',
        height: 240,
        offsetY: -30,
        offsetX: 20,
        dropShadow: {
            enabled: true,
            top: 5,
            left: 0,
            bottom: 0,
            right: 0,
            blur: 5,
            color: '#b6c2e4',
            opacity: 0.35
        },
    },
    plotOptions: {
        radialBar: {
            inverseOrder: true,
            hollow: {
                margin: 5,
                size: '55%',
                background: 'transparent',
            },
            track: {
                show: true,
                background: '#ddd',
                strokeWidth: '10%',
                opacity: 1,
                margin: 5, // margin is in pixels
            },

            dataLabels: {
                name: {
                    fontSize: '18px',
                },
                value: {
                    fontSize: '16px',
                    color: '#50649c',
                },
            }
        },
    },
    series: [0],
    labels: ['Usage'],
    legend: {
        show: true,
        position: 'bottom',
        offsetX: -40,
        offsetY: -10,
        formatter: function (val, opts) {
            return val + " - " + opts.w.globals.series[opts.seriesIndex] + '%'
        }
    },
    fill: {
        type: 'gradient',
        gradient: {
            shade: 'dark',
            type: 'horizontal',
            shadeIntensity: 0.5,
            inverseColors: true,
            opacityFrom: 1,
            opacityTo: 1,
            stops: [0, 100],
            gradientToColors: ["#ffb822", "#5d78ff", "#34bfa3"]
        }
    },
    stroke: {
        lineCap: 'round'
    },
}
var optionsCircle3 = {
    chart: {
        type: 'radialBar',
        height: 240,
        offsetY: -30,
        offsetX: 20,
        dropShadow: {
            enabled: true,
            top: 5,
            left: 0,
            bottom: 0,
            right: 0,
            blur: 5,
            color: '#b6c2e4',
            opacity: 0.35
        },
    },
    plotOptions: {
        radialBar: {
            inverseOrder: true,
            hollow: {
                margin: 5,
                size: '55%',
                background: 'transparent',
            },
            track: {
                show: true,
                background: '#ddd',
                strokeWidth: '10%',
                opacity: 1,
                margin: 5, // margin is in pixels
            },

            dataLabels: {
                name: {
                    fontSize: '18px',
                },
                value: {
                    fontSize: '16px',
                    color: '#50649c',
                },
            }
        },
    },
    series: [0, 0, 0],
    labels: ['Total', 'Available', 'Used'],
    legend: {
        show: true,
        position: 'bottom',
        offsetX: -40,
        offsetY: -10,
        formatter: function (val, opts) {
            return val + " - " + opts.w.globals.series[opts.seriesIndex] + '%'
        }
    },
    fill: {
        type: 'gradient',
        gradient: {
            shade: 'dark',
            type: 'horizontal',
            shadeIntensity: 0.5,
            inverseColors: true,
            opacityFrom: 1,
            opacityTo: 1,
            stops: [0, 100],
            gradientToColors: ["#ffb822", "#5d78ff", "#34bfa3"]
        }
    },
    stroke: {
        lineCap: 'round'
    },
}

try {
    var chartCircle = new ApexCharts(document.querySelector('#circlechart'), optionsCircle);
    chartCircle.render();
    var circlechartMempru = new ApexCharts(document.querySelector('#circlechartMempru'), optionsCircle2);
    circlechartMempru.render();
    var memoryUses = new ApexCharts(document.querySelector('#ana_device'), optionsCircle3);
    memoryUses.render();
} catch (e) {

}


var iteration = 11

function getRandom() {
    var i = iteration;
    return (Math.sin(i / trigoStrength) * (i / trigoStrength) + i / trigoStrength + 1) * (trigoStrength * 2)
}

function getRangeRandom(yrange) {
    return Math.floor(Math.random() * (yrange.max - yrange.min + 1)) + yrange.min
}

//
// window.setInterval(function () {
//
//     iteration++;
//
//     // $.ajax({
//     //     url: "admin/cpuinfo",
//     //     type: 'GET',
//     //     context: this,
//     //     success: function (response, textStatus, jQxhr) {
//     //         const events = [];
//     //         chartCircle.updateSeries([ response])
//     //     },
//     //     error: function (jqXhr, textStatus, errorThrown) {
//     //
//     //         console.log(errorThrown);
//     //     }
//     // });
//
//     // $.ajax({
//     //     url: "admin/mem-cpu",
//     //     type: 'GET',
//     //     context: this,
//     //     success: function (response, textStatus, jQxhr) {
//     //         const events = [];
//     //         console.log(response.usedPercent)
//     //         chartCircle.updateSeries([response.cpu])
//     //         circlechartMempru.updateSeries([parseInt(response.memory)])
//     //     },
//     //     error: function (jqXhr, textStatus, errorThrown) {
//     //
//     //         console.log(errorThrown);
//     //     }
//     // });
//
//
// }, 3000)


var randomizeArray = function (arg) {
    var array = arg.slice();
    var currentIndex = array.length, temporaryValue, randomIndex;

    while (0 !== currentIndex) {

        randomIndex = Math.floor(Math.random() * currentIndex);
        currentIndex -= 1;

        temporaryValue = array[currentIndex];
        array[currentIndex] = array[randomIndex];
        array[randomIndex] = temporaryValue;
    }

    return array;
}

// data for the sparklines that appear below header area
var sparklineData = [47, 45, 54, 38, 56, 24, 65, 31, 37, 39, 62, 51, 35, 41, 35, 27, 93, 53, 61, 27, 54, 43, 19, 46];


var dash_spark_1 = {

    chart: {
        type: 'area',
        height: 85,
        sparkline: {
            enabled: true
        },
        dropShadow: {
            enabled: true,
            top: 12,
            left: 0,
            bottom: 5,
            right: 0,
            blur: 2,
            color: '#8997bd',
            opacity: 0.1
        },
    },
    stroke: {
        curve: 'smooth',
        width: 3
    },
    fill: {
        opacity: 1,
        gradient: {
            shade: '#2c77f4',
            type: "horizontal",
            shadeIntensity: 0.5,
            inverseColors: true,
            opacityFrom: 0.1,
            opacityTo: 0.1,
            stops: [0, 80, 100],
            colorStops: []
        },
    },
    series: [{
        data: [4, 8, 5, 10, 4, 16, 5, 11, 6, 11, 30, 10, 13, 4, 6, 3, 6]
    }],
    colors: ['#2c77f4'],
}
new ApexCharts(document.querySelector("#dash_spark_1"), dash_spark_1).render();
var data = [{
    data: []
}]

var spark2 = {
    chart: {
        id: 'realtime',
        height: 350,
        type: 'line',
        animations: {
            enabled: true,
            easing: 'linear',
            dynamicAnimation: {
                speed: 1000
            }
        },
        dropShadow: {
            enabled: true,
            top: 10,
            left: 0,
            bottom: 0,
            right: 0,
            blur: 2,
            color: '#b6c2e4',
            opacity: 0.35
        },
    },

    stroke: {
        width: 2,
        curve: 'smooth'
    },
    fill: {
        opacity: 0.2,
    },
    series: data,
    yaxis: {
        min: 0
    },
    colors: ['#fbb624'],
    title: {
        text: 'Trafice ;ive',
        offsetX: 20,
        style: {
            fontSize: '24px'
        }
    },
    subtitle: {
        text: 'Expenses',
        offsetX: 20,
        style: {
            fontSize: '14px'
        }
    }
}
var lastDate = 0;


traficChart = document.querySelector("#trafics")// ApexCharts(document.querySelector("#spark2"), dash_spark_1);
// traficChart.render()

var conn;
var arraydata = []
if (window["WebSocket"]) {
    conn = new WebSocket("ws://" + document.location.host + "/ws");
    conn.onopen = function () {
        console.log("<p>Socket is open</p>");
    };
    conn.onclose = function (evt) {
        var item = document.createElement("div");
        item.innerHTML = "<b>Connection closed.</b>";
        appendLog(item);
    };
    conn.onmessage = function (evt) {
        // var messages = evt.data.split('\n');
        try {
            const result = JSON.parse(evt.data)
            if (window.location.pathname === "/admin/" || window.location.pathname === "/admin") {
                if (chartCircle) {
                    chartCircle.updateSeries([result.cpu])
                }

                if (circlechartMempru) {
                    circlechartMempru.updateSeries([parseInt(result.memory.usedPercent)])
                }

                if (memoryUses) {
                    const used = (result.memory.used * 100) / result.memory.total
                    const available = (result.memory.available * 100) / result.memory.total
                    memoryUses.updateSeries([100, parseInt(available), parseInt(used),])

                }
            }
            let sendByte = 0
            console.log(window.location.pathname)
            if (window.location.pathname === "/admin/" || window.location.pathname === "/admin" || window.location.pathname === "/admin/network" || window.location.pathname === "/admin/network/") {
                row = ""
                for (let i = 0; i < result.netinfo.length; i++) {
                    row += "      <tr>\n" +
                        "                    <th>" + result.netinfo[i].name + "</th>\n" +
                        "                    <th>" + bytesToSize(result.netinfo[i].bytesSent) + "</th>\n" +
                        "                    <th>" + bytesToSize(result.netinfo[i].bytesRecv) + "</th>\n" +
                        "                </tr>"

                    // arraydata.push()
                }

                traficChart.innerHTML = row

            }

        } catch (e) {
            console.log(e)
        }


    };
} else {

}
const units = ['bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

function bytesToSize(x) {
    let l = 0, n = parseInt(x / 1024, 10) || 0;

    while (n >= 1024 && ++l) {
        n = n / 1024;
    }

    return (n.toFixed(n < 10 && l > 0 ? 1 : 0) + ' ' + units[l]);
}