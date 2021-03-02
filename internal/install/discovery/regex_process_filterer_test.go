// +build unit

package discovery

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestShouldFilterWithCmdLineInsteadOfName(t *testing.T) {
	r := []types.Recipe{
		{
			ID:           "1",
			Name:         "test-cassandra-ohi",
			ProcessMatch: []string{"cassandra"},
		},
	}

	processes := []types.GenericProcess{
		mockProcess{
			name:    "java",
			cmdline: "java -xyz processSomething/cassandra",
		},
		mockProcess{
			name:    "somethingElse",
			cmdline: "somethingElse",
		},
	}

	mockRecipeFetcher := recipes.NewMockRecipeFetcher()
	mockRecipeFetcher.FetchRecipesVal = r
	f := NewRegexProcessFilterer(mockRecipeFetcher)
	filtered, err := f.filter(context.Background(), processes, types.DiscoveryManifest{})

	require.NoError(t, err)
	require.NotNil(t, filtered)
	require.NotEmpty(t, filtered)
	require.Equal(t, 1, len(filtered))
	require.Equal(t, filtered[0].MatchingPattern, "cassandra")
}

func TestFilter_NoMatchingProcess(t *testing.T) {
	r := []types.Recipe{
		{
			ID:           "1",
			Name:         "java-agent",
			ProcessMatch: []string{"java"},
		},
		{
			ID:           "2",
			Name:         "cassandra-open-source-integration",
			ProcessMatch: []string{"cassandra", "cassandradaemon", "cqlsh"},
		},
		{
			ID:           "3",
			Name:         "jmx-open-source-integration",
			ProcessMatch: []string{"java.*jboss", "java.*tomcat", "java.*jetty"},
		},
	}

	processes := []types.GenericProcess{
		mockProcess{
			name:    "nonMatchingProcess",
			cmdline: "nonMatchingProcess",
		},
	}

	mockRecipeFetcher := recipes.NewMockRecipeFetcher()
	mockRecipeFetcher.FetchRecipesVal = r
	f := NewRegexProcessFilterer(mockRecipeFetcher)
	filtered, err := f.filter(context.Background(), processes, types.DiscoveryManifest{})

	require.NoError(t, err)
	require.NotNil(t, filtered)
	require.Equal(t, 0, len(filtered))
}

func TestFilter_SingleMatchingProcess_SingleOHIRecipe(t *testing.T) {
	r := []types.Recipe{
		{
			ID:           "1",
			Name:         "cassandra-open-source-integration",
			ProcessMatch: []string{"cassandra", "cassandradaemon", "cqlsh"},
		},
	}

	processes := []types.GenericProcess{
		mockProcess{
			name:    "cassandra",
			cmdline: "/usr/lib/jvm/java-1.8.0-openjdk-1.8.0.272.b10-1.amzn2.0.1.x86_64/jre/bin/java -Xloggc:/var/log/cassandra/gc.log -ea -XX:+UseThreadPriorities -XX:ThreadPriorityPolicy=42 -XX:+HeapDumpOnOutOfMemoryError -Xss256k -XX:StringTableSize=1000003 -XX:+AlwaysPreTouch -XX:-UseBiasedLocking -XX:+UseTLAB -XX:+ResizeTLAB -XX:+UseNUMA -XX:+PerfDisableSharedMem -Djava.net.preferIPv4Stack=true -XX:+UseParNewGC -XX:+UseConcMarkSweepGC -XX:+CMSParallelRemarkEnabled -XX:SurvivorRatio=8 -XX:MaxTenuringThreshold=1 -XX:CMSInitiatingOccupancyFraction=75 -XX:+UseCMSInitiatingOccupancyOnly -XX:CMSWaitDuration=10000 -XX:+CMSParallelInitialMarkEnabled -XX:+CMSEdenChunksRecordAlways -XX:+CMSClassUnloadingEnabled -XX:+PrintGCDetails -XX:+PrintGCDateStamps -XX:+PrintHeapAtGC -XX:+PrintTenuringDistribution -XX:+PrintGCApplicationStoppedTime -XX:+PrintPromotionFailure -XX:+UseGCLogFileRotation -XX:NumberOfGCLogFiles=10 -XX:GCLogFileSize=10M -Xms977M -Xmx977M -Xmn200M -XX:+UseCondCardMark -XX:CompileCommandFile=/etc/cassandra/conf/hotspot_compiler -javaagent:/usr/share/cassandra/lib/jamm-0.3.0.jar -Dcassandra.jmx.local.port=7199 -Dcom.sun.management.jmxremote.authenticate=false -Dcom.sun.management.jmxremote.password.file=/etc/cassandra/jmxremote.password -Djava.library.path=/usr/share/cassandra/lib/sigar-bin -XX:OnOutOfMemoryError=kill -9 %p -Dlogback.configurationFile=logback.xml -Dcassandra.logdir=/var/log/cassandra -Dcassandra.storagedir= -Dcassandra-pidfile=/var/run/cassandra/cassandra.pid -cp /etc/cassandra/conf:/usr/share/cassandra/lib/airline-0.6.jar:/usr/share/cassandra/lib/antlr-runtime-3.5.2.jar:/usr/share/cassandra/lib/asm-5.0.4.jar:/usr/share/cassandra/lib/caffeine-2.2.6.jar:/usr/share/cassandra/lib/cassandra-driver-core-3.0.1-shaded.jar:/usr/share/cassandra/lib/commons-cli-1.1.jar:/usr/share/cassandra/lib/commons-codec-1.9.jar:/usr/share/cassandra/lib/commons-lang3-3.1.jar:/usr/share/cassandra/lib/commons-math3-3.2.jar:/usr/share/cassandra/lib/compress-lzf-0.8.4.jar:/usr/share/cassandra/lib/concurrentlinkedhashmap-lru-1.4.jar:/usr/share/cassandra/lib/concurrent-trees-2.4.0.jar:/usr/share/cassandra/lib/disruptor-3.0.1.jar:/usr/share/cassandra/lib/ecj-4.4.2.jar:/usr/share/cassandra/lib/guava-18.0.jar:/usr/share/cassandra/lib/HdrHistogram-2.1.9.jar:/usr/share/cassandra/lib/high-scale-lib-1.0.6.jar:/usr/share/cassandra/lib/hppc-0.5.4.jar:/usr/share/cassandra/lib/jackson-annotations-2.9.10.jar:/usr/share/cassandra/lib/jackson-core-2.9.10.jar:/usr/share/cassandra/lib/jackson-databind-2.9.10.4.jar:/usr/share/cassandra/lib/jamm-0.3.0.jar:/usr/share/cassandra/lib/javax.inject.jar:/usr/share/cassandra/lib/jbcrypt-0.3m.jar:/usr/share/cassandra/lib/jcl-over-slf4j-1.7.7.jar:/usr/share/cassandra/lib/jctools-core-1.2.1.jar:/usr/share/cassandra/lib/jflex-1.6.0.jar:/usr/share/cassandra/lib/jna-4.2.2.jar:/usr/share/cassandra/lib/joda-time-2.4.jar:/usr/share/cassandra/lib/json-simple-1.1.jar:/usr/share/cassandra/lib/jstackjunit-0.0.1.jar:/usr/share/cassandra/lib/libthrift-0.9.2.jar:/usr/share/cassandra/lib/log4j-over-slf4j-1.7.7.jar:/usr/share/cassandra/lib/logback-classic-1.1.3.jar:/usr/share/cassandra/lib/logback-core-1.1.3.jar:/usr/share/cassandra/lib/lz4-1.3.0.jar:/usr/share/cassandra/lib/metrics-core-3.1.5.jar:/usr/share/cassandra/lib/metrics-jvm-3.1.5.jar:/usr/share/cassandra/lib/metrics-logback-3.1.5.jar:/usr/share/cassandra/lib/netty-all-4.0.44.Final.jar:/usr/share/cassandra/lib/ohc-core-0.4.4.jar:/usr/share/cassandra/lib/ohc-core-j8-0.4.4.jar:/usr/share/cassandra/lib/reporter-config3-3.0.3.jar:/usr/share/cassandra/lib/reporter-config-base-3.0.3.jar:/usr/share/cassandra/lib/sigar-1.6.4.jar:/usr/share/cassandra/lib/slf4j-api-1.7.7.jar:/usr/share/cassandra/lib/snakeyaml-1.11.jar:/usr/share/cassandra/lib/snappy-java-1.1.1.7.jar:/usr/share/cassandra/lib/snowball-stemmer-1.3.0.581.1.jar:/usr/share/cassandra/lib/ST4-4.0.8.jar:/usr/share/cassandra/lib/stream-2.5.2.jar:/usr/share/cassandra/lib/thrift-server-0.3.7.jar:/usr/share/cassandra/apache-cassandra-3.11.10.jar:/usr/share/cassandra/apache-cassandra-thrift-3.11.10.jar:/usr/share/cassandra/stress.jar: org.apache.cassandra.service.CassandraDaemon",
		},
		mockProcess{
			name:    "somethingElse",
			cmdline: "somethingElse",
		},
	}

	mockRecipeFetcher := recipes.NewMockRecipeFetcher()
	mockRecipeFetcher.FetchRecipesVal = r
	f := NewRegexProcessFilterer(mockRecipeFetcher)
	filtered, err := f.filter(context.Background(), processes, types.DiscoveryManifest{})

	require.NoError(t, err)
	require.NotNil(t, filtered)
	require.NotEmpty(t, filtered)
	require.Equal(t, 1, len(filtered))
	require.Equal(t, filtered[0].MatchingPattern, "cassandra")
}

func TestFilter_SingleMatchingProcess_SingleAPMRecipe(t *testing.T) {
	r := []types.Recipe{
		{
			ID:           "1",
			Name:         "java-agent",
			ProcessMatch: []string{"java"},
		},
	}

	processes := []types.GenericProcess{
		mockProcess{
			name:    "java",
			cmdline: "java -D[Standalone] -server -Xms64m -Xmx512m -XX:MetaspaceSize=96M -XX:MaxMetaspaceSize=256m -Djava.net.preferIPv4Stack=true -Djboss.modules.system.pkgs=org.jboss.byteman -Djava.awt.headless=true --add-exports=java.base/sun.nio.ch=ALL-UNNAMED --add-exports=jdk.unsupported/sun.misc=ALL-UNNAMED --add-exports=jdk.unsupported/sun.reflect=ALL-UNNAMED -Dorg.jboss.boot.log.file=/opt/wildfly/standalone/log/server.log -Dlogging.configuration=file:/opt/wildfly/standalone/configuration/logging.properties -jar /opt/wildfly/jboss-modules.jar -mp /opt/wildfly/modules org.jboss.as.standalone -Djboss.home.dir=/opt/wildfly -Djboss.server.base.dir=/opt/wildfly/standalone -c standalone.xml -b 0.0.0.0",
		},
		mockProcess{
			name:    "somethingElse",
			cmdline: "somethingElse",
		},
	}

	mockRecipeFetcher := recipes.NewMockRecipeFetcher()
	mockRecipeFetcher.FetchRecipesVal = r
	f := NewRegexProcessFilterer(mockRecipeFetcher)
	filtered, err := f.filter(context.Background(), processes, types.DiscoveryManifest{})

	require.NoError(t, err)
	require.NotNil(t, filtered)
	require.NotEmpty(t, filtered)
	require.Equal(t, 1, len(filtered))
	require.Equal(t, filtered[0].MatchingPattern, "java")
}

func TestFilter_SingleMatchingProcess_MultipleRecipes(t *testing.T) {
	r := []types.Recipe{
		{
			ID:           "1",
			Name:         "java-agent",
			ProcessMatch: []string{"java"},
		},
		{
			ID:           "2",
			Name:         "cassandra-open-source-integration",
			ProcessMatch: []string{"cassandra", "cassandradaemon", "cqlsh"},
		},
	}

	processes := []types.GenericProcess{
		mockProcess{
			name:    "cassandra",
			cmdline: "/usr/lib/jvm/java-1.8.0-openjdk-1.8.0.272.b10-1.amzn2.0.1.x86_64/jre/bin/java -Xloggc:/var/log/cassandra/gc.log -ea -XX:+UseThreadPriorities -XX:ThreadPriorityPolicy=42 -XX:+HeapDumpOnOutOfMemoryError -Xss256k -XX:StringTableSize=1000003 -XX:+AlwaysPreTouch -XX:-UseBiasedLocking -XX:+UseTLAB -XX:+ResizeTLAB -XX:+UseNUMA -XX:+PerfDisableSharedMem -Djava.net.preferIPv4Stack=true -XX:+UseParNewGC -XX:+UseConcMarkSweepGC -XX:+CMSParallelRemarkEnabled -XX:SurvivorRatio=8 -XX:MaxTenuringThreshold=1 -XX:CMSInitiatingOccupancyFraction=75 -XX:+UseCMSInitiatingOccupancyOnly -XX:CMSWaitDuration=10000 -XX:+CMSParallelInitialMarkEnabled -XX:+CMSEdenChunksRecordAlways -XX:+CMSClassUnloadingEnabled -XX:+PrintGCDetails -XX:+PrintGCDateStamps -XX:+PrintHeapAtGC -XX:+PrintTenuringDistribution -XX:+PrintGCApplicationStoppedTime -XX:+PrintPromotionFailure -XX:+UseGCLogFileRotation -XX:NumberOfGCLogFiles=10 -XX:GCLogFileSize=10M -Xms977M -Xmx977M -Xmn200M -XX:+UseCondCardMark -XX:CompileCommandFile=/etc/cassandra/conf/hotspot_compiler -javaagent:/usr/share/cassandra/lib/jamm-0.3.0.jar -Dcassandra.jmx.local.port=7199 -Dcom.sun.management.jmxremote.authenticate=false -Dcom.sun.management.jmxremote.password.file=/etc/cassandra/jmxremote.password -Djava.library.path=/usr/share/cassandra/lib/sigar-bin -XX:OnOutOfMemoryError=kill -9 %p -Dlogback.configurationFile=logback.xml -Dcassandra.logdir=/var/log/cassandra -Dcassandra.storagedir= -Dcassandra-pidfile=/var/run/cassandra/cassandra.pid -cp /etc/cassandra/conf:/usr/share/cassandra/lib/airline-0.6.jar:/usr/share/cassandra/lib/antlr-runtime-3.5.2.jar:/usr/share/cassandra/lib/asm-5.0.4.jar:/usr/share/cassandra/lib/caffeine-2.2.6.jar:/usr/share/cassandra/lib/cassandra-driver-core-3.0.1-shaded.jar:/usr/share/cassandra/lib/commons-cli-1.1.jar:/usr/share/cassandra/lib/commons-codec-1.9.jar:/usr/share/cassandra/lib/commons-lang3-3.1.jar:/usr/share/cassandra/lib/commons-math3-3.2.jar:/usr/share/cassandra/lib/compress-lzf-0.8.4.jar:/usr/share/cassandra/lib/concurrentlinkedhashmap-lru-1.4.jar:/usr/share/cassandra/lib/concurrent-trees-2.4.0.jar:/usr/share/cassandra/lib/disruptor-3.0.1.jar:/usr/share/cassandra/lib/ecj-4.4.2.jar:/usr/share/cassandra/lib/guava-18.0.jar:/usr/share/cassandra/lib/HdrHistogram-2.1.9.jar:/usr/share/cassandra/lib/high-scale-lib-1.0.6.jar:/usr/share/cassandra/lib/hppc-0.5.4.jar:/usr/share/cassandra/lib/jackson-annotations-2.9.10.jar:/usr/share/cassandra/lib/jackson-core-2.9.10.jar:/usr/share/cassandra/lib/jackson-databind-2.9.10.4.jar:/usr/share/cassandra/lib/jamm-0.3.0.jar:/usr/share/cassandra/lib/javax.inject.jar:/usr/share/cassandra/lib/jbcrypt-0.3m.jar:/usr/share/cassandra/lib/jcl-over-slf4j-1.7.7.jar:/usr/share/cassandra/lib/jctools-core-1.2.1.jar:/usr/share/cassandra/lib/jflex-1.6.0.jar:/usr/share/cassandra/lib/jna-4.2.2.jar:/usr/share/cassandra/lib/joda-time-2.4.jar:/usr/share/cassandra/lib/json-simple-1.1.jar:/usr/share/cassandra/lib/jstackjunit-0.0.1.jar:/usr/share/cassandra/lib/libthrift-0.9.2.jar:/usr/share/cassandra/lib/log4j-over-slf4j-1.7.7.jar:/usr/share/cassandra/lib/logback-classic-1.1.3.jar:/usr/share/cassandra/lib/logback-core-1.1.3.jar:/usr/share/cassandra/lib/lz4-1.3.0.jar:/usr/share/cassandra/lib/metrics-core-3.1.5.jar:/usr/share/cassandra/lib/metrics-jvm-3.1.5.jar:/usr/share/cassandra/lib/metrics-logback-3.1.5.jar:/usr/share/cassandra/lib/netty-all-4.0.44.Final.jar:/usr/share/cassandra/lib/ohc-core-0.4.4.jar:/usr/share/cassandra/lib/ohc-core-j8-0.4.4.jar:/usr/share/cassandra/lib/reporter-config3-3.0.3.jar:/usr/share/cassandra/lib/reporter-config-base-3.0.3.jar:/usr/share/cassandra/lib/sigar-1.6.4.jar:/usr/share/cassandra/lib/slf4j-api-1.7.7.jar:/usr/share/cassandra/lib/snakeyaml-1.11.jar:/usr/share/cassandra/lib/snappy-java-1.1.1.7.jar:/usr/share/cassandra/lib/snowball-stemmer-1.3.0.581.1.jar:/usr/share/cassandra/lib/ST4-4.0.8.jar:/usr/share/cassandra/lib/stream-2.5.2.jar:/usr/share/cassandra/lib/thrift-server-0.3.7.jar:/usr/share/cassandra/apache-cassandra-3.11.10.jar:/usr/share/cassandra/apache-cassandra-thrift-3.11.10.jar:/usr/share/cassandra/stress.jar: org.apache.cassandra.service.CassandraDaemon",
		},
		mockProcess{
			name:    "somethingElse",
			cmdline: "somethingElse",
		},
	}

	mockRecipeFetcher := recipes.NewMockRecipeFetcher()
	mockRecipeFetcher.FetchRecipesVal = r
	f := NewRegexProcessFilterer(mockRecipeFetcher)
	filtered, err := f.filter(context.Background(), processes, types.DiscoveryManifest{})

	require.NoError(t, err)
	require.NotNil(t, filtered)
	require.NotEmpty(t, filtered)
	require.Equal(t, 2, len(filtered))
	require.Equal(t, filtered[0].MatchingPattern, "java")
	require.Equal(t, filtered[1].MatchingPattern, "cassandra")
}

func TestFilter_MultipleMatchingProcesses_SingleRecipe(t *testing.T) {
	r := []types.Recipe{
		{
			ID:           "1",
			Name:         "test-java-agent",
			ProcessMatch: []string{"java"},
		},
	}

	processes := []types.GenericProcess{
		mockProcess{
			name:    "cassandra",
			cmdline: "java -xyz processSomething/cassandra",
		},
		mockProcess{
			name:    "java",
			cmdline: "java -D[Standalone] -server -Xms64m -Xmx512m -XX:MetaspaceSize=96M -XX:MaxMetaspaceSize=256m -Djava.net.preferIPv4Stack=true -Djboss.modules.system.pkgs=org.jboss.byteman -Djava.awt.headless=true --add-exports=java.base/sun.nio.ch=ALL-UNNAMED --add-exports=jdk.unsupported/sun.misc=ALL-UNNAMED --add-exports=jdk.unsupported/sun.reflect=ALL-UNNAMED -Dorg.jboss.boot.log.file=/opt/wildfly/standalone/log/server.log -Dlogging.configuration=file:/opt/wildfly/standalone/configuration/logging.properties -jar /opt/wildfly/jboss-modules.jar -mp /opt/wildfly/modules org.jboss.as.standalone -Djboss.home.dir=/opt/wildfly -Djboss.server.base.dir=/opt/wildfly/standalone -c standalone.xml -b 0.0.0.0",
		},
		mockProcess{
			name:    "somethingElse",
			cmdline: "somethingElse",
		},
	}

	mockRecipeFetcher := recipes.NewMockRecipeFetcher()
	mockRecipeFetcher.FetchRecipesVal = r
	f := NewRegexProcessFilterer(mockRecipeFetcher)
	filtered, err := f.filter(context.Background(), processes, types.DiscoveryManifest{})

	require.NoError(t, err)
	require.NotNil(t, filtered)
	require.NotEmpty(t, filtered)
	require.Equal(t, 2, len(filtered))
	require.Equal(t, filtered[0].MatchingPattern, "java")
	require.Equal(t, filtered[1].MatchingPattern, "java")
}

func TestFilter_MultipleMatchingProcesses_MultipleRecipes(t *testing.T) {
	r := []types.Recipe{
		{
			ID:           "1",
			Name:         "test-java-agent",
			ProcessMatch: []string{"java"},
		},
		{
			ID:           "2",
			Name:         "jmx-open-source-integration",
			ProcessMatch: []string{"java.*jboss", "java.*tomcat", "java.*jetty"},
		},
		{
			ID:           "3",
			Name:         "cassandra-open-source-integration",
			ProcessMatch: []string{"cassandra", "cassandradaemon", "cqlsh"},
		},
	}

	processes := []types.GenericProcess{
		mockProcess{
			name:    "cassandra",
			cmdline: "java -xyz processSomething/cassandra",
		},
		mockProcess{
			name:    "java",
			cmdline: "java -D[Standalone] -server -Xms64m -Xmx512m -XX:MetaspaceSize=96M -XX:MaxMetaspaceSize=256m -Djava.net.preferIPv4Stack=true -Djboss.modules.system.pkgs=org.jboss.byteman -Djava.awt.headless=true --add-exports=java.base/sun.nio.ch=ALL-UNNAMED --add-exports=jdk.unsupported/sun.misc=ALL-UNNAMED --add-exports=jdk.unsupported/sun.reflect=ALL-UNNAMED -Dorg.jboss.boot.log.file=/opt/wildfly/standalone/log/server.log -Dlogging.configuration=file:/opt/wildfly/standalone/configuration/logging.properties -jar /opt/wildfly/jboss-modules.jar -mp /opt/wildfly/modules org.jboss.as.standalone -Djboss.home.dir=/opt/wildfly -Djboss.server.base.dir=/opt/wildfly/standalone -c standalone.xml -b 0.0.0.0",
		},
		mockProcess{
			name:    "somethingElse",
			cmdline: "somethingElse",
		},
	}

	mockRecipeFetcher := recipes.NewMockRecipeFetcher()
	mockRecipeFetcher.FetchRecipesVal = r
	f := NewRegexProcessFilterer(mockRecipeFetcher)
	filtered, err := f.filter(context.Background(), processes, types.DiscoveryManifest{})

	require.NoError(t, err)
	require.NotNil(t, filtered)
	require.NotEmpty(t, filtered)
	require.Equal(t, 4, len(filtered))
	require.Equal(t, filtered[0].MatchingPattern, "java")
	require.Equal(t, filtered[1].MatchingPattern, "cassandra")
	require.Equal(t, filtered[2].MatchingPattern, "java")
	require.Equal(t, filtered[3].MatchingPattern, "java.*jboss")
}
