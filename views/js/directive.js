/**
 * Created by rang on 2015-06-07.
 */
var DIRECTIVE_URL = '/views/js/template/directives';
angular.module('common.directives', []).
    directive('navbar', [function() {
        return  {
            templateUrl : DIRECTIVE_URL + "/navbar.html",
            restrict : "AE",
            transclude: true,
            scope: {
                totalCount: '=',
                currentPage: '=',
                pageSize: '@',
                pagingColSize: '@',
                selectedCallBack: '&'
            },
            link: function (scope, el, attr) {

            }
        };
    }]);