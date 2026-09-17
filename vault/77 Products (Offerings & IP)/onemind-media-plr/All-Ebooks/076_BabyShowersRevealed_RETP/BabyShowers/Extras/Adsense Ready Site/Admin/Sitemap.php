<?php include("Session.php"); ?>
	<?php include("Header.php"); ?>
		<center>
			<table class="tblForm">
				<tr>
					<td align="center" style="background: url(http://www.niche-mania.com/Images/NM_006.jpg); color: #FFFFFF; font-weight: bold; height: 24px; text-align: center; vertical-align: middle;">
						XML Sitemap Generator &nbsp; 
					</td>				
				</tr>
				<?php
					if(!$_SESSION['XMLSitemap'])
						{
				?>
							<tr>
								<td>
									<p align="justify">
										After configuring or modifying your Niche Mania Website, the last step is to build or rebuild the 
										XML Sitemap that you will submit to Google for thorough indexing of your site. 
										<br>
										<br>
										To do this, simply click the "Generate XML Sitemap" button below - the XML Sitemap will automatically be generated for you.
									</p>
									<br>
									<img align="left" src="http://www.niche-mania.com/Images/NM_004.jpg" hspace="5" vspace="5">
									<p align="justify">
										<small>
											<b>
												CAUTION:
											</b>
											Although this process normally only takes a few seconds, DO NOT navigate away from this page until it has completed. 
											Doing so will corrupt your XML Sitemap and Google will not be able to read it!
										</small>
									</p>
								</td>				
							</tr>
								<tr>
									<td>
										<hr size="0">
									</td>				
								<tr>
							<tr>
								<td class="tdFieldValue">
									<input class="iButton" onClick="this.disabled='true'; location.href='XML.php';" type="button" value=" Generate XML Sitemap ">
								</td>				
							</tr>
				<?php
						}
					else
						{
							$_SESSION['XMLSitemap'] = false
				?>
							<tr>
								<td>
									<p align="justify">
										Congratulations - your XML Sitemap has been generated and saved in your root directory. 
										You must now submit it to Google Sitemaps via the following link:
										<br>
										<br>
										<center>
											<a href="http://www.google.com/webmasters/sitemaps" target="_new">http://www.google.com/webmasters/sitemaps</a>
											<br>
											<br>
											URL to submit to Google Sitemaps (copy and paste it):
											<br>
											<br>
											<b>
												<?php echo $_SESSION['Domain']."Sitemap.xml"; ?>
											</b>										
										</center>
										<br>
										
										
										NOTE: This system requires that you have a valid Google Gmail account. If you need one, please send email to:
										<br>
										<br>
										<center>
											<a href="mailto:michalis.kotzakolios@gmail.com">michalis.kotzakolios@gmail.com</a>
										</center>
										<br>
										Be certain to include the email address where you would like us to send your Gmail Invitation.										
										<br>
										<br>																				
									</p>
								</td>				
							</tr>
								<tr>
									<td>
										<hr size="0">
									</td>				
								<tr>
							<tr>
								<td class="tdFieldValue">
									<input class="iButton" onClick="this.disabled='true'; window.open('../Sitemap.xml');" type="button" value=" Preview XML Sitemap ">
								</td>				
							</tr>				
				<?php
						}				
				?>
			</table>
		</center>
	<?php include("Footer.php"); ?>