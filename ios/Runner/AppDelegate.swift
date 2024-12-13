import Flutter
import UIKit
import flutter_local_notifications
import workmanager

@main
@objc class AppDelegate: FlutterAppDelegate {
  override func application(
    _ application: UIApplication,
    didFinishLaunchingWithOptions launchOptions: [UIApplication.LaunchOptionsKey: Any]?
  ) -> Bool {
    //GeneratedPluginRegistrant.register(with: self)
      

    FlutterLocalNotificationsPlugin.setPluginRegistrantCallback { (registry) in
          GeneratedPluginRegistrant.register(with: registry)
     }

     

     if #available(iOS 10.0, *) {
       UNUserNotificationCenter.current().delegate = self as UNUserNotificationCenterDelegate
     }

      
     WorkmanagerPlugin.setPluginRegistrantCallback { registry in
                 GeneratedPluginRegistrant.register(with: registry)
             }
      


     WorkmanagerPlugin.registerBGProcessingTask(withIdentifier: "com.example.beehive.rescheduledTask")
    WorkmanagerPlugin.registerPeriodicTask(withIdentifier: "com.example.beehive.simplePeriodicTask", frequency: NSNumber(value: 15 * 60))
      UIApplication.shared.setMinimumBackgroundFetchInterval(TimeInterval(60*15))

     WorkmanagerPlugin.registerPeriodicTask(withIdentifier: "com.example.beehive.iOSBackgroundAppRefresh", frequency: NSNumber(value: 15 * 60))



    GeneratedPluginRegistrant.register(with: self)
    return super.application(application, didFinishLaunchingWithOptions: launchOptions)

  }
}
