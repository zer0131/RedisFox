/* Info Widget
* ====================== */

var InfoWidget = BaseWidget.extend({

    initialize: function () {

        this.Name = "Info Widget";

        this.init();
        this.updateFrequency = 30000;

        // templates
        var templateSource = $("#info-widget-template").html()
            , popOverTemplateSource = $("#popover-template").html()
            , infoTemplateSource = $("#info-template").html();

        this.template = Handlebars.compile(templateSource);
        this.popOverTemplate = Handlebars.compile(popOverTemplateSource);
        this.infoTemplate = Handlebars.compile(infoTemplateSource);

    },

    render: function () {

        var model = this.model.toJSON()
            , markUp = this.template(model)
            , popoverMarkup = this.popOverTemplate(model.databases)
            , infoMarkup = this.infoTemplate(model);

        $(this.el).html(markUp);

        $('#total-keys').popover({"title": "数据详情", "content": popoverMarkup});

        $('#misc-info').popover({"title": "Info", "content": infoMarkup, "placement": "bottom"});
    }

});