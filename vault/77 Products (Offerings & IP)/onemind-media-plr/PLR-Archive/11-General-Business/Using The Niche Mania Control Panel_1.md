---
id: plr-20260114013944-1287
title: "Using The Niche Mania Control Panel 1"
type: reference
path: PLR
domain: reference
status: reference
created: 2026-01-14
modified: 2026-01-14
owner: zeus
tags: [path/plr, type/reference, status/reference, topic/general-business]
ai_summary: "PLR reference - Using The Niche Mania Control Panel 1"
---
# Using The Niche Mania Control Panel_1

      Niche Mania Control Panel

              USER GUIDE
                            Version 3.0



                Michalis “BIG Mike” Kotzakolios
 Content Mania | Niche Mania | Niche Maniacs | Xpro Mania




© 2006, Michalis “BIG Mike” Kotzakolios – All Rights Reserved Worldwide
Niche Mania Control Panel User's Guide                                                                          Page 2 of 16



Table of Contents
      INTRODUCTION.............................................................................................................. 3

      SYSTEM REQUIREMENTS................................................................................................... 3

      INSTALLATION STRATEGY................................................................................................. 3

      INSTALLATION............................................................................................................... 4

      FIRST LOGIN................................................................................................................. 4

      CONTROL PANEL NAVIGATION........................................................................................... 4

      LOAD DEFAULTS............................................................................................................ 5

      SITE CONFIGURATION...................................................................................................... 5

      ADSENSE SETTINGS......................................................................................................... 8

      YPN SETTINGS............................................................................................................... 8

      AMAZON SETTINGS......................................................................................................... 9

      ADMIN SECURITY............................................................................................................ 9

      CONTACT INFORMATION.................................................................................................. 9

      LINK PARTNERS............................................................................................................ 10

      RSS CONFIGURATION..................................................................................................... 10

      ARTICLE MANAGER........................................................................................................ 10

      PUBLISH WEBSITE......................................................................................................... 12

      LOG OUT.................................................................................................................... 13

      ADDITIONAL STRATEGIES................................................................................................ 13

      OTHER PRODUCTS AND SERVICES FROM NICHE MANIA............................................................. 16




   © 2006, Michalis “BIG Mike” Kotzakolios – All Rights Reserved Worldwide…
Niche Mania Control Panel User's Guide                             Page 3 of 16


INTRODUCTION

The Niche Mania Control Panel is a platform dependent Administrative System
for managing our Niche Mania Adsense Content Sites. It allows you to
completely administer your sites through a web browser, from anywhere in the
world you may happen to be.

The Control Panel is a work in progress and continues to evolve as the
marketing techniques and ideas we encounter change over time. Our goal is to
provide an intuitive, user friendly system that is productive to use for both
beginners and experts alike.

If you have a suggestion for the Control Panel, you can email it to
mike@anamik.com – we take all suggestions seriously and implement those that
we feel will benefit the broadest range of users.

Note: Bear in mind that, by sending us your suggestions, you are giving us the
right to use them, royalty free, in our Control Panel design.

Unless you have some expertise with JavaScript and PHP, we recommend that
you do not attempt to edit your Niche Mania Adsense Sites in a typical HTML
Editor – these often corrupt scripting or add superfluous HTML that interferes
with the overall look, feel and operation of the sites.


SYSTEM REQUIREMENTS

   ●   PHP 4.3 or Higher
   ●   Ability to set File and Folder Permissions


INSTALLATION STRATEGY

Niche Mania Adsense Content Sites may be installed in the root of a domain, a
folder, sub-folder or on a sub-domain or a folder of a sub-domain. The
underlying scripting will automatically recognize and interact with the type of
installation you performed.

You may mix and match different types of installations on the domain – that is
entirely up to you and your own approach to Adsense. We'll discuss additional
strategies at the end of the User Guide to help you in planning future
installations.




    © 2006, Michalis “BIG Mike” Kotzakolios – All Rights Reserved Worldwide…
Niche Mania Control Panel User's Guide                              Page 4 of 16


INSTALLATION

After you've downloaded and unzipped your Niche Mania Adsense Content Site,
you will need to upload it using an FTP Client. You should upload it while
preserving the folder and file structure originally provided. Do not rename or
move any of the folders or files – doing so will cause unexpected problems with
the site's operation.

After uploading your site, you will also need to set a few File and Folder
permissions as indicated below:

   ●   RssFeed.xml, Sitemap.xml, urllist.txt – set to 646
   ●   Aticles and Content Folders – set to 747
   ●   All Files in the Content and Includes Folders – set to 646

Note: If the Control Panel does not work properly after setting the permissions
above, try resetting all of the Folder and File permissions to 755 – some hosts
require this.

Note: do not change any of the other permissions set by default.


FIRST LOGIN

After completing the installation, you may then navigate to your site through a
web browser to the location you installed. As this was a new installation, you
will automatically be redirected to the Control Panel Login page.

The default User Name is “admin” and the Password is “password”. After
logging in, you'll receive a message to change these and will automatically
redirect to the Admin Security page to do so.


CONTROL PANEL NAVIGATION

Navigating the Control Panel is as simple as clicking any button in the left-hand
column. Each button leads to a page to configure and/or work with the portion
of your site indicated.

You can navigate to, and work with any section of the Control Panel as desired.
It does not require you to add information in any particular order. Note
however, that if you fail to save your changes before navigating to a new page,
any changes made on the current page will be lost.



    © 2006, Michalis “BIG Mike” Kotzakolios – All Rights Reserved Worldwide…
Niche Mania Control Panel User's Guide                                Page 5 of 16


LOAD DEFAULTS

If you have previously entered this domain and configured it in the Niche Mania
Portal, you may remotely load your default settings by first clicking the Load
Defaults button and then clicking the Save Default Configuration button on that
page.


When you do this, the Control Panel contacts the Portal; looks up the domain
name and then looks for any Default Configuration you may have set up. This
includes:


   ●   Adsense Publisher's ID
   ●   Amazon Associate ID
   ●   Contact Informational
   ●   Yahoo Partner's ID


It then loads and saves this information, which is a huge time saver when
installing multiple sites.


SITE CONFIGURATION

This page is used to define the “Look and Feel” of your site, how Adsense and
YPN is implemented and also the Keywords you may be targeting.


The first field, Domain Name is disabled and is displayed only for informational
purposes. It displays the root URL that it believes it is installed on. This is also
the URL that you would use for promoting this site, such as when doing link
exchanges or article submissions.


The Site Name field can be confusing – whatever you enter here is displayed on
a horizontal menu bar prefixed to the menu items and in the footer of every
page. A simple title is all that is required – you do not need to “Stuff” it with
keywords as the site is already optimized.


Keep it simple, i.e.; Acne Information or just Acne!




    © 2006, Michalis “BIG Mike” Kotzakolios – All Rights Reserved Worldwide…
Niche Mania Control Panel User's Guide                              Page 6 of 16


The Keywords field is pre-populated with the default Keywords the site is
targeting. You may add, edit or delete these keywords to suite your own needs.
Make sure that each Keyword or phrase is followed by a comma and a space –
except for the last Keyword in the list.


When you add a Keyword, it becomes available in the Article Manager to assign
to your articles and other content. Since your site also comes with one article
for each Keyword, it is important that you not delete the existing Keywords
without also deleting the articles associated with them.


Both the Main Title and Sub-Title fields are by far the least understood in our
Adsense sites. Your Adsense site was delivered with a professionally designed
header image that blends well with the site. It also included a blank header
image that does not contain any text and is intended to be used in conjunction
with these two fields.


To do so, you need to rename the blank header image to H0.jpg and upload it
to the Images folder. After doing so, you may enter an appropriate Main Title
and Sub-title that will be displayed over the blank header image.


Because of the variety of designs used, you may find it necessary to edit the
Styles.css file located in the Includes Folder to position the text from these
fields to suite your needs. You're going to look for and carefully edit this
section of the CSS:


h1
       {
             color:         #000000;
             font-size:     small;
             margin-left:   5;
             margin-top:    20;
       }

h3
       {
             color:         #000000;
             font-size:     small;
             margin-left:   5;
             margin-top:    20;
       }


By changing the values for the highlighted items, you can position the text
anywhere over the blank header image – then save and upload this file.


     © 2006, Michalis “BIG Mike” Kotzakolios – All Rights Reserved Worldwide…
Niche Mania Control Panel User's Guide                             Page 7 of 16


The PPC Settings radio buttons allow you to choose which PPC Program to
target for this site. If you are an Adsense Publisher you would of course leave
the default of Adsense selected. Conversely, if you are a Yahoo Partner, you
would select that instead.


If you are in both PPC programs, you may select Random, to randomly display
one or the other program's ad units (they will never be displayed together).
Over a 24 hour period, you will typically see a fairly even distribution of ad
units for both programs.


If you have not yet applied for either program, you can choose None, which
will prevent any ad units from being displayed. This is especially useful if you
are using the site when applying for your Adsense account.


The Rotate Articles check box, determines how the articles are displayed to
visitors who navigate to your site without specifying a page name. Each time a
visitor navigates to the domain itself (without indicating a particular page), a
random article is displayed. This is useful in having Search Engines more
thoroughly crawl and index the pages on your site.


If you prefer to have a static home page, you can un-tick this check box, which
in turn enables the Home Page drop-down list. You may choose any article to
act as your static Home Page – if you prefer to have a custom Home Page, do
the following:


   1. Navigate to the Article Manager and click the Add Article button.
   2. Using the built-in HTML Editor, create enter the content for your new
      Home Page.
   3. After saving your changes, navigate back to the Site Configuration page
      and select your new Home Page from the drop-down list.


After doing so, whenever a visitor arrives on your site without specifying a
particular page, they will always see your static Home Page rather than a
random article.


The last configurable item on this page is the Layout drop-down list. You can
choose from one of three Layouts to display your site. All three layouts are
highly optimized for SEO and PPC so your choice is really one of personal
preference.


    © 2006, Michalis “BIG Mike” Kotzakolios – All Rights Reserved Worldwide…
Niche Mania Control Panel User's Guide                             Page 8 of 16


ADSENSE SETTINGS

If you are an Adsense Publisher, you will need to enter your Adsense
Publisher's ID in the first field like this: pub-6754415291760887


Note the prefixed “pub-” and that there are no extra spaces at the end of the
number. If you do not enter it exactly as displayed, you will not receive credit
for any clicks on your ads.


The next field is for entering an optional Channel Code. Channel Codes allow
you to do some very basic tracking of your Adsense earnings by assigning a
Name and Number to for tracking purposes. You will always enter a number in
this field as assigned by Google.


Rather than attempt to explain how to use Channel Codes, we recommend that
you view Google's help files online regarding these.


The Link and URL Color fields are pre-populated with a matching color scheme
to allow your ad units to blend into the overall site design. The Link Color
controls the color of the title of each Adsense ad, while the URL Color controls
the color of the URL displayed at the bottom of each ad, if any.


You may want to experiment with these colors to identify the best performing
combination for your site. The colors we selected will normally perform well,
but even minor changes may serve to increase your CTR.


Lastly for this page, you may tick or un-tick any of the various options to
display a variety of Google Referral Buttons and optional ad units. There are
pros and cons to each item, so it's best to experiment to determine what
converts best for the topic of your site.


YPN SETTINGS

This page works exactly like the Adsense Settings page – fill in the Partner's ID
field with your ID, enter an optional Reporting Category (Channel Code) and
save your settings.




    © 2006, Michalis “BIG Mike” Kotzakolios – All Rights Reserved Worldwide…
Niche Mania Control Panel User's Guide                             Page 9 of 16


AMAZON SETTINGS

Another way to monetize a Niche Mania Adsense Content site is through the
display of Amazon ads. This page is pre-populated with most settings – you only
need to enter your Associate's ID and tick the Display Amazon check box to
enable them.


You can refine your ad targeting by selecting a Category from the drop-down
list. Each Category is pre-defined by Amazon and they are broad enough to
cover virtually any topic.


The Keyword field is pre-populated with a primary Keyword to target;
however, if you decide to change Keywords, you need to be aware of one
important item when using multiple words or phrases. Each space between
words should be replaced with the + (Plus) symbol. This tells Amazon to
concatenate the words when searching for ad to display.


ADMIN SECURITY

This page is used to change your User Name and Password. Please note that
you may not use the default combination of “admin” and “password” and that
until you do change them, a “Nag” dialog is displayed as a reminder.


Passwords are encrypted and stored in the Configure.php file located in the
Includes Folder. If you ever forget your User Name or Password, simply
download the Configure.php file, open it in Windows Notepad and carefully
change the variables $UserName and $Password (the text between the
adjacent quotes). Then upload the file and login.


CONTACT INFORMATION

All Contact Information you enter on this page is displayed on the Legal
Information page of your site. This information is optional and not required.
You may enter as much or as little Contact Information as desired to suite your
needs.


Note: Minimally, you should consider entering a Name and Email Address.


    © 2006, Michalis “BIG Mike” Kotzakolios – All Rights Reserved Worldwide…
Niche Mania Control Panel User's Guide                            Page 10 of 16


LINK PARTNERS

The Link Partners page is used to enter reciprocal links that you may be
exchanging with other sites or to link to your own sites around the net. It uses
a built-in HTML Editor, to allow you to “Dress Up” your links and descriptions.


RSS CONFIGURATION

Niche Mania Adsense Content Sites use two types of RSS Feeds for News and
Links. These are pre-configured for you to deliver relevant news and links from
Yahoo and MSN. You may, at your discretion, replace these feeds with your
own favorites. As long as they conform to RSS 2.0, they will work normally.


You may also specify how many items are displayed from the feeds by editing
the values in the News Quantity and Links Quantity fields. We recommend
that you keep these numbers small as the larger the number, the longer a given
page will take to load.


The Ping Trigger field determines how many page views must occur between
automatically pinging directories to indicate that some content has changed
(such as a news Feed update). Since you rarely want to ping any directory more
than once a day, set this number to the average number of unique visitors your
receive per day.


You may also “Turn Off” News or Links Feeds by un-ticking the corresponding
check box. We recommend against this as the feeds automatically provide
fresh content for the Search Engines.


ARTICLE MANAGER

The Article Manager uses an online HTML Editor to simplify the creation and
editing of your articles.


When you click the Manage Articles button, the page reloads and displays a
drop-down list containing all of your articles and two button (Add Articles and
Load Article).



    © 2006, Michalis “BIG Mike” Kotzakolios – All Rights Reserved Worldwide…
Niche Mania Control Panel User's Guide                              Page 11 of 16


To Add a new article, simply click the Add Article button and the page will
reload with the HTML Editor ready for you to write your article!


The Keyword field is a drop-down list containing all of the Keywords you're
targeting with this site. In order to assign this article to a Keyword, select one
from the list. If you are targeting a new Keyword, you should have previously
entered it on the Site Configuration page.


Next, you can schedule when your articles will first appear by entering a future
date in the Publication Date field. This field allows you to prepare several
articles at once and Publish or “Drip” them over time.


Note that for best results, the date should be in this format: mm/dd/yyyy such
as 07/11/2006 for July 7, 2006.


As part of this feature, a page will automatically be created as a placeholder to
allow Search Engines to “Pre-Index” the article, pending it's actual publication.


The Article Title field is critical to how the entire system works – it should be
in Title Case (first letter of each word capitalized), each word separated by
spaces and should NOT contain any punctuation.


NOTE: When you save this article, it will be saved by replacing the spaces
between the words with underscores and appending a PHP extension to it. For
example, if your article title is “My First Article On Dog Treats”, it will then be
converted to “My_First_Article_On_Dog_Treats.php”


When you subsequently edit this article, even if you change the Title, the file
name will not change. The reason for this is because as your articles are
indexed by the Search Engines, you do not want to lose a listing because you
arbitrarily changed the Title.


Lastly for the Title field, you should incorporate the Keyword into the Title,
but it cannot be only the Title. This is because a separate article page is
created for each Keyword to act as a mini-index for all articles targeting those
Keywords.


If you use a Keyword as a Title, the article will never be displayed!


    © 2006, Michalis “BIG Mike” Kotzakolios – All Rights Reserved Worldwide…
Niche Mania Control Panel User's Guide                            Page 12 of 16


You can enter your own name or a pseudonym in the Author's Name field. The
choice is entirely up to you.


You will actually type or paste your article into the HTML Editor, which works
in much the same way as any typical word processing program. You have a
variety of formatting controls on the Tool Bar and should spend a few minutes
familiarizing yourself with their operation.


There is a second HTML Editor at the bottom of the page – this used to create
an optional Resource Box (About the Author) for your article. If you use it, this
information will be displayed at the bottom of the article.


When finished, click the Save Article button to save the article, or click the
Cancel button to abandon your article and return to the Article Manager.


If you want to Edit and existing article, you must first select it from the drop-
down list in the Article Manager. After doing so, click the Load Article button.
Once you've Loaded the article, you can either Delete it (there is a
confirmation screen to confirm doing so) or Edit it by clicking the appropriate
button.


Editing an article works exactly like Adding an article – all of the fields and
controls are the same.


PUBLISH WEBSITE

The Niche Mania Control Panel uses an Article Library, which stores the raw
text of each article in the Content Folder of the site. When you Add, Edit or
Delete an article, you are not actually doing anything to the PHP files in the
Articles Folder. All changes occur in Content Folder until you Publish your
Website.


Whenever you make any changes, you should always perform one last action
before logging out – clicking the Publish Website button. Doing so rebuilds your
website, incorporating any changes and also updates your XML Sitemap for
Google, the URL List for Yahoo and the RSS Feed for your own site (the feed
that visitors will subscribe to). It also pings the major directories to announce
your changes in real time!



    © 2006, Michalis “BIG Mike” Kotzakolios – All Rights Reserved Worldwide…
Niche Mania Control Panel User's Guide                              Page 13 of 16


LOG OUT

If you are hosting multiple Niche Mania Adsense Content Sites, especially in
folders, you must be certain to use the Log Out button when exiting the
Control Panel. Doing so clears variables that are stored in memory – if you fail
to do so, log into another Niche Mania site and begin to work, you may
inadvertently overwrite that site's settings with the previous one!


ADDITIONAL STRATEGIES

There are practically an infinite number of ways to deploy your Niche Mania
Adsense sites, and I'll share my preferences with you.


A Super Niche Site is my favorite as it's easier to maintain. In this strategy, you
simply register a generic domain name, and install each site either in Folders
or Sub-Domains. On the root of the domain you can install a blog or a Portal
Landing Page that links to all of the individual sites.


NOTE: You can cross-link related sites in this strategy, but should be very
careful to avoid cross-linking unrelated sites, which may be viewed as spam.


A Broadly Targeted Niche Site is my second favorite – in this scenario, you set
up 12 to 16 domains following a typical Search Engine Directory structure and
then install related sites on each domain, cross-linking everything together.
Again, use blogs or Portal Landing Pages in the root of the domain.


A Highly Focused Niche Site is my final preference, although these take more
work. To do this, you choose a niche, set up a domain and concentrate just on
that niche and nothing else.


In the order listed, a Super Niche Site is the least costly and easiest solution to
maintain, while a Highly Focused Niche Site is the most costly and difficult to
maintain. Earnings wise, they're all similar – you're trading off time and effort
for volume, which works well.




    © 2006, Michalis “BIG Mike” Kotzakolios – All Rights Reserved Worldwide…
Niche Mania Control Panel User's Guide                           Page 14 of 16


      Are You Using Social Bookmarking
   As Part of your Traffic Building Strategy?

 Wouldn't it Be Great If There Was A Way To
 Automate The Process? Well Guess What?

                    SOCIALBOT
  A Social Bookmarking Program for Windows From Niche Mania




  I Want To Show You How To Quickly And
 Easily Bookmark All Of Your Website's Pages

          (Without Doing It Page-By-Page)

    With SocialBot - The HOT New Social
  Bookmarking Program From Niche Mania!


   © 2006, Michalis “BIG Mike” Kotzakolios – All Rights Reserved Worldwide…
Niche Mania Control Panel User's Guide                           Page 15 of 16



      Are You Sick And Tired Of All The
  "Get Rich With Adsense" Programs On The
                  Internet?

        Would You Like Hard Core Data,
   Statistically Analyzed By An Engineer and
          Recognized Adsense Expert?

      The Truth About Adsense Profitability

     A Special Report on the Real Earnings Potential for New Webmasters




    I Want To Show You What Other People,
   Just Like Yourself Are Really Earning and
            How They Are Doing It!


   © 2006, Michalis “BIG Mike” Kotzakolios – All Rights Reserved Worldwide…
Niche Mania Control Panel User's Guide                            Page 16 of 16


OTHER PRODUCTS AND SERVICES FROM NICHE MANIA

Niche Mania primarily focuses on products and services related to PPC
Publishing and Content Style websites. We presently offer twelve membership
choices spread over three primary markets:


                                   Content Mania – This membership group is
                                   designed to deliver PLR Articles, Graphics,
                                   Reports and other PLR products to it's
                                   members. More information may be found
                                   here: http://www.content-mania.com



Niche Mania – The “Original” Mania
Group, is designed to deliver high quality,
easy management Adsense Content Sites
to it's members. In addition, the Niche
Mania Portal provides additional tools
integrated with our Adsense Content Sites
to make them by far the most profitable
on the Internet.


Niche Maniacs Forum – The forum is a meeting place for members to learn,
share and brainstorm new ideas and techniques on how to market their content
sites. With over a thousand members, it's a fun and useful forum to regularly
visit at http://forum.my-webtools.net


                                              Xpro Mania – This member ship
                                              group focuses on those customers
                                              who use XSite Pro and delivers
                                              high quality sites and templates to
                                              it's   members.        For    more
                                              information please go here:
                                              http://www.xpro-mania.com




No matter which Mania group you prefer, you can rest assured that you will
always receive unparalleled support and guidance from our staff of experts!



    © 2006, Michalis “BIG Mike” Kotzakolios – All Rights Reserved Worldwide…
