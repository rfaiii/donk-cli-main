import SwiftUI

@main
struct DonkCompanionApp: App {
    @StateObject private var connectionStore = ConnectionStore()
    @StateObject private var sessionStore = SessionStore()

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environmentObject(connectionStore)
                .environmentObject(sessionStore)
                .preferredColorScheme(.dark)
                .tint(Color.primary)
        }
    }
}
