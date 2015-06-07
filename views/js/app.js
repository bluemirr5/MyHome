/**
 * Created by rang on 2015-06-07.
 */
var PARTIAL_URL = '/views/partial';
angular.module('myHome', [
    'ngRoute',
    'common.directives',
    'remotejob.controllers'
])
    .config(['$routeProvider', function($routeProvider) {
        $routeProvider.when('/remotejob', {	templateUrl: PARTIAL_URL+'/remotejob/remote-job.html',	controller: 'RemoteJobController'	});
        $routeProvider.when('/remotemeeting',{	templateUrl: PARTIAL_URL+'/remotemeeting/remote-meeting.html',	controller: ''	});
        $routeProvider.when('/contact', {	templateUrl: PARTIAL_URL+'/contact/contact.html',	controller: ''	});
        $routeProvider.otherwise({	redirectTo: '/remotejob'	});
    }]);
