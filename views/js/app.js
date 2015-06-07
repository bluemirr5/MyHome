/**
 * Created by rang on 2015-06-07.
 */
var PARTIAL_URL = '/views/partial';
angular.module('myHome', [
    'ngRoute',
    'common.directives',
    'remotejob.controllers',
    'remotemeeting.controllers',
])
.config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/remotejob', {	templateUrl: PARTIAL_URL+'/remote-job.html',	controller: 'RemoteJobController'	});
    $routeProvider.when('/remotemeeting',{	templateUrl: PARTIAL_URL+'/remote-meeting.html',	controller: 'RemoteMeetingController'	});
    $routeProvider.when('/contact', {	templateUrl: PARTIAL_URL+'/contact.html',	controller: ''	});
    $routeProvider.otherwise({	redirectTo: '/remotejob'	});
}]);
