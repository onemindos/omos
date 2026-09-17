// Bookmark this page
// copyright 12th February 2005, 21st June 2006
// Stephen Chapman, http://javascript.about.com
// Updated for Opera 9
// You have permission to copy and use this javascript provided that
// the copyright notice is included with the script

var chr = 'CTRL-D';
var agt=navigator.userAgent.toLowerCase();
if(agt.substr(agt.indexOf('opera')+6,1) < 9) chr = 'CTRL-T';
if (window.external) document.write('<a href="javascript:window.external.AddFavorite(self.location,document.title)">Add to Favorites<\/a>');
else document.write('Press '+chr+' to bookmark this page.');
                  