/**
 * Created by rang on 2015-06-07.
 */
var DIRECTIVE_URL = '/views/js/directives';
angular.module('common.directives', []).
    directive('navbar', ['$location', function($location) {
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
                scope.changePage = function(pageId) {
                    $location.url(pageId);
                };
            }
        };
    }]);