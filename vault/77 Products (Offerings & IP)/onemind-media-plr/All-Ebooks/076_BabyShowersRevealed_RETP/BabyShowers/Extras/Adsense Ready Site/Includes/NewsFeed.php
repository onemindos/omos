<?php

  error_reporting(0);
	
  $CAT  = $_GET['CAT'];
  $CAT = str_replace(" ", "+", $CAT);
  $QTY = 10; //$_GET['QTY'];
  
  $url  = "http://news.search.yahoo.com/news/rss?p=<keywords>";
  $url = str_replace("<keywords>", $CAT, $url);

	print "<ul class='ulCategory'>";
	flush();

class SimpleXmlParser{

	var $SimpleParser;
	var $feedUrl;

	var $title = "";
	var $description = "";
	var $link = "";
	var $author="";
	var $pubDate="";
    var $insideitem = false;
	var $tag = "";
	var $count_results = 0;
	var $total_results = 0;
	
	function SimpleXmlParser($MyFeed, $results) {
		
		$this->total_results = $results;		
	    $this->SimpleParser = xml_parser_create();           
		$this->feedUrl=$MyFeed;   
		
		xml_set_object($this->SimpleParser,$this);	
		xml_set_element_handler($this->SimpleParser, "XmlParserFirstElement", "XmlParserendElement");
		xml_set_character_data_handler($this->SimpleParser, "characterData");  
		
		$this->ParseFeed();   // Call to Parser Function
	}
	
	function XmlParserFirstElement($parser, $tagName, $attrs) {

		if ($this->insideitem) {
			$this->tag = $tagName;
		} elseif ($tagName == "ITEM") {
			$this->insideitem = true;
		}
	}

    function XmlParserendElement($parser, $tagName) {
        //The Function Will be called, when ever the XML_PARSER Encounters a end Tag, in the XML File
		if ($tagName == "ITEM") {

			if ($this->count_results < $this->total_results) {

			  print "<li>";
			  print "  <a target=_new href=\"http://" . $this->link . "\">";
			  print      htmlspecialchars(trim($this->title));
			  print "  </a><br><br>";
			  print htmlspecialchars(trim($this->description)) . "<br>";			
			  print "</li><br>";								
						
			  $this->title = "";
			  $this->description = "";
			  $this->link = "";
			  $this->insideitem = false;
			}
			$this->count_results = $this->count_results + 1;
		}
	}	
	
    function characterData($parser, $data) {

		if ($this->insideitem) {
			switch ($this->tag) {
				case "TITLE":
				$this->title .= $data;
				break;
				case "DESCRIPTION":
				$this->description .= $data;
				break;

				case "LINK":
				  $this->link .= $data;
				  $pos = strpos($this->link, "*");
				  if (!($pos === false)) {				  
				    $this->link  = substr($this->link , $pos + 1, 300);					  
				  }
				  $this->link = trim($this->link);
				  $pos = strpos($this->link, "//");
				  if (!($pos === false)) {				  
				    $this->link  = substr($this->link , $pos + 2, 300);					  
				  }				  
				  break;
			}
		}
	}

	function ParseFeed(){

		$data = file_get_contents($this->feedUrl);
		xml_parse($this->SimpleParser, $data);		
        xml_parser_free($this->SimpleParser);        
	}
}	

  $XMLpar = new SimpleXmlParser($url, $QTY);
  
  print "</ul>"
  
?>
