import {
  moduleForModel,
  test
} from 'ember-qunit';

moduleForModel('branch', {
  // Specify the other units that are required for this test.
  needs: ['model:project', 'model:build']
});

test('it exists', function(assert) {
  var model = this.subject();
  // var store = this.store();
  assert.ok(!!model);
});
