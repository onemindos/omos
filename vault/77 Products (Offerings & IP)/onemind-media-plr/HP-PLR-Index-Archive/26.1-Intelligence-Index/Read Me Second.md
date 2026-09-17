---
id: hp-26-20260113225955
title: "Read Me Second"
type: reference
path: HP
domain: 25-49
folder: "26"
status: active
created: 2026-01-13
modified: 2026-01-13
owner: zeus
tags: [path/hp, path/hp/26, type/reference, status/active]
ai_summary: "Read Me Second"
---
# Read Me Second

First let me say that writing directions is not my thing. So if at any point I lose you please feel free to contact me. I will try to make myself as available to everyone as possible. Right now you can get a hold of me by posting to the forum, PM me through the forum, and/or e-mail me at alan.sean@gmail.com. I am trying to also get the time to set up a Skype account and will let you know my user name there as soon as I do. I am going to keep these instructions very simple and go step by step. 

You will need a Search and Replace program of some kind to do a mass search/replace for your user IDs. What I use and recommend is a program called "Simple Search/Replace" this is an open source program that is included in your AdSense sites package. If you have another program that you are comfortable using feel free to use that instead. However, these instructions will be based on using "Simple Search/Replace" and I may not be able to answer questions if you are using another program.

After unzipping the site you should have 1 folder, AdSense site, and 4 files "Read Me First.doc", "search_setup.exe", Home-Template.php, Article-Template.php. As you are aware this is the
"Read First.doc " 

Install "Simple Search/Replace"

Open "Simple Search/Replace" in the Path to search: field browse to the folder containing the AdSense site.

Make sure the box for Search Subdirectories is checked.

In the pull down menu for Select file types: select *.php if that is not an option then simply type in *.php (* is shift+8).

In the Find what text: field type in the search phrase (all search phrases are listed at the end of this PDF) and in the Replace with: put your new information. Make sure you have a check in the box for Case Sensitive field.

Push the Start button.

Verify the information is correct and press YES. If it is not correct press NO.

If you have not selected to make a back-up file it will tell you that there is no back-up file and ask if you still want to continue. Press Yes. I personally don't use the back-up option as I have a back-up copy of the site in the original zip file.

It will then show you all pages that the changes were made on.

Click the Search/Replace tab.

Repeat steps 5 through 10 till you have replaced all the search terms.

Once you have replaced all the search terms your site is ready to upload to your server.

Search Phrases that you will need to change:

YourGoogleID = Your Google AdSense ID Number ONLY. This is the number directly after pub-

YourAmazonID = Your Amazon Associate's ID

YourClickBankID = Your ClickBank user name

YourPayDotComID = Your PayDotCom user name

YourDomainName = The name of your domain and/or site name. This is used on the home page to welcome users to your site so be sure to enter the name you want your site known as.

ScriptLoco = The location on your site that you will be placing the scripts. Do Not include http:// I recommend that this be your route directory. So if your installing the site to it's own domain or sub-domain (which I would strongly recommend as some of the features don't work if put in a folder/directory) you would simply replace this with my-domain.com/ or sub-domain.my-domain.com/. If you are setting the site up in a folder (not recommended) then it would be my-domain.com/my-directory/.


Editing Your Web Site:

I am including an article template and an index template. These are both php files but don't let that bother you. You can open them in any HTML editor. Note: The template is going to look strange as it will not have direct access to the style sheet and images. This is ok as it will when you upload it back to the sever. 

For the Home-Template.php replace "Your content goes here" with your content. Do not change anything else or the site might not work. Then upload the changed template to your server. What I do and recommend you do as well is then rename the index.php file on your sever to index-old.php and then change Home-Template.php  to index.php. This way if there is an error on the new index you can quickly put the old one back up. I highly recommend that you at least create your own home page to help keep your site original.

The Article-Template.php has three areas for you to modifie.

At the very top of the article table you will see "PUT YOUR KEYWORD HERE" this is where you will put the article keyword
Next is "TITLE OF YOUR ARTICLE" replace this with the new article title.
Last Your content goes here is where you place the content of your new article.

Though creating a new article page is simple enough there are currently 2 main complications. 

The first being that the new article is not automatcly added to any of the navigation menus. Therefore it is not linked to your site in anyway. I am working on this. However for the time being what I recommend you do is to have a section of your modified home page list new articles and link to any new articles you put up in that fashion. The other thing you can do is take a snippet of the article and place that snippet into the randomlink.txt file. Just fallow the formatting in the file. By placing your snippet in that file it will rotate through your pages. If you want the snippet seen on all pages then you have to place it on all of the randomlink.txt files (there is 1 file per directory).

The second one is that your new articles will not have the Master Keyword hyperlinked to an affiliate product.
