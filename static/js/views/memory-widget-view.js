var MemoryWidget = BaseWidget.extend({

    initialize: function () {

        this.Name = "Memory Widget";

        this.init();
        this.updateFrequency = 30000;

        // templates
        var templateSelector = "#memory-widget-template"
            , templateSource = $(templateSelector).html();

        this.template = Handlebars.compile(templateSource);
        this.$el.empty().html(this.template());

        // chart
        this.chart = new google.visualization.LineChart($("#memory-widget-chart").empty().get(0));
        this.dataTable = new google.visualization.DataTable();
        this.dataTable.addColumn('datetime', 'DateTime');
        this.dataTable.addColumn('number', 'Max');
        this.dataTable.addColumn('number', 'Current');
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
                var timeArr = datetimeArr[1].split(":");
                var recordDate = new Date(parseInt(dateArr[0]), parseInt(dateArr[1]) - 1, parseInt(dateArr[2]), parseInt(timeArr[0]), parseInt(timeArr[1]), parseInt(timeArr[2]));

                if (self.dataTable)
                    self.dataTable.addRow([recordDate, parseInt(obj[1]), parseInt(obj[2])]);
            });

            var pointSize = model.data.length > 120 ? 1 : 5, options = {
                title: '',
                colors: ['#1581AA', '#77BA44'],
                pointSize: pointSize,
                chartArea: {'top': 10, 'width': '85%'},
                width: "100%",
                height: 200,
                animation: {duration: 500, easing: 'out'}
            };

            this.chart.draw(this.dataTable, options);
        }
    }
});