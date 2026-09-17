									</td>
								</tr>
							</table>
						</td>
					</tr>	
					<tr>							
						<td class="tdFooterInfo">					
							<?php 
									echo $_SERVER["REMOTE_ADDR"]." | ";
									
									if (isset($_SESSION['UserName']))
										{
											echo $_SERVER["HTTP_HOST"]." | ".$_SERVER["DOCUMENT_ROOT"];
										}
									else
										{
											echo "Welcome to the Niche Mania Website Control Panel...";
										} 											
							?>
						</td>
					</tr>
					<tr>				
						<td class="tdFooter">					
							<p>
								Powered By Niche Mania | &copy 2005 - <a href="http://www.niche-mania.com" target="_new">Kotzakoliou, SSA</a> | All Rights Reserved Worldwide 
							</p>
						</td>
					</tr>
				</table>				
		</center>
	</body>
</html>