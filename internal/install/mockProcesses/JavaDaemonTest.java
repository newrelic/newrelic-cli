public class JavaDaemonTest {
    public static void main(String args[]) {
        Runnable daemonRunner = new Runnable() {
            public void run() {
                while (true) {
                    try {
                        Thread.sleep(500);
                    } catch (InterruptedException ignored) {
                    }
                }
            }
        };
        Thread daemonThread = new Thread(daemonRunner);
        daemonThread.setDaemon(true);
        daemonThread.start();
        try {
            Thread.sleep(30000);
        } catch (InterruptedException ignored) {
        }
        System.out.println("Done.");
    }
}