/*
wifibytes is a package for calculating the number of kilobits per
seconds are going over a given network interface.

   client, err := wifibytes.NewClient(10, "wlp2s0")
   if err != nil {
	   // handle error
   }
   down, up, err := client.GetKilobitsPerSecond()
   if err != nil {
	   // handle error
   }
   // We now have download speed and upload speed
*/
package wifibytes
