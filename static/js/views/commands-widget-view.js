var CommandsWidget = BaseWidget.extend({

    initialize: function () {

        this.Name = "Commands Widget";

        this.init();
        this.updateFrequency = 30000;

        // templates
        var templateSelector = "#commands-widget-template"
            , templateSource = $(templateSelector).html();

        this.template = Handlebars.compile(templateSource);
        this.$el.empty().html(this.template());

        // chart
        this.chart = new google.visualization.AreaChart($("#commands-widget-chart").empty().get(0));
        this.dataTable = new google.visualization.DataTable();
        this.dataTable.addColumn('datetime', 'DateTime');
        this.dataTable.addColumn('number', 'Commands Processed');
    },

    render: function () {

        var model = this.model.toJSON(), markUp = this.template(model), self = this;

        self.dataTable.removeRows(0, self.dataTable.getNumberOfRows());

        if (model.data != null) {
            $.each(model.data, function (index, obj) {

                // first item of the object contains datetime info
                // [ YYYY, MM, DD, HH, MM, SS ]
                var datetimeArr = obj[0].split(" ");
                var dateArr = datetimeArr[0].split("-");
                var timeArr;
                if (datetimeArr.length > 1) {
                    timeArr = datetimeArr[1].split(":");
                    if (timeArr.length == 1) {
                        timeArr[1] = 0;
                        timeArr[2] =0
                    } else if (timeArr.length == 2) {
                        timeArr[2] = 0;
                    }
                } else {
                    timeArr = [0, 0, 0];
                }
                var recordDate = new Date(parseInt(dateArr[0]), parseInt(dateArr[1]) - 1, parseInt(dateArr[2]), parseInt(timeArr[0]), parseInt(timeArr[1]), parseInt(timeArr[2]));

                if (self.dataTable)
                    self.dataTable.addRow([recordDate, parseInt(obj[1])]);
            });


            var pointSize = model.data.length > 120 ? 1 : 5, options = {
                title: '',
                colors: ['#17BECF', '#9EDAE5'],
                areaOpacity: .9,
                pointSize: pointSize,
                chartArea: {'top': 10, 'width': '85%'},
                width: "100%",
                height: 200,
                animation: {duration: 500, easing: 'out'},
                vAxis: {minValue: 0}
            };

            this.chart.draw(this.dataTable, options);
        }
    }
});